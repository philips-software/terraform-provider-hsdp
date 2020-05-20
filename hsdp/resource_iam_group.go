package hsdp

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
	"net/http"
)

func resourceIAMGroup() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceIAMGroupCreate,
		Read:   resourceIAMGroupRead,
		Update: resourceIAMGroupUpdate,
		Delete: resourceIAMGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"managing_organization": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"roles": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"users": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 2000,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"services": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 2000,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceIAMGroupCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	var group iam.Group
	group.Description = d.Get("description").(string)
	group.Name = d.Get("name").(string)
	group.ManagingOrganization = d.Get("managing_organization").(string)

	createdGroup, _, err := client.Groups.CreateGroup(group)
	if err != nil {
		return err
	}
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	d.SetId(createdGroup.ID)
	d.Set("name", createdGroup.Name)
	d.Set("description", createdGroup.Description)
	d.Set("managing_organization", createdGroup.ManagingOrganization)

	// Add roles
	for _, r := range roles {
		role, _, _ := client.Roles.GetRoleByID(r)
		if role != nil {
			client.Groups.AssignRole(*createdGroup, *role)
		}
	}

	// Add users
	users := expandStringList(d.Get("users").(*schema.Set).List())
	if len(users) > 0 {
		_, _, err = client.Groups.AddMembers(*createdGroup, users...)
	}

	// Add services
	services := expandStringList(d.Get("services").(*schema.Set).List())
	if len(services) > 0 {
		_, _, err = client.Groups.AddServices(*createdGroup, services...)
	}
	return err
}

func resourceIAMGroupRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	id := d.Id()
	group, resp, err := client.Groups.GetGroupByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	_ = d.Set("managing_organization", group.ManagingOrganization)
	_ = d.Set("description", group.Description)
	_ = d.Set("name", group.Name)
	roles, _, err := client.Groups.GetRoles(*group)
	if err != nil {
		return err
	}
	roleIDs := make([]string, len(*roles))
	for i, r := range *roles {
		roleIDs[i] = r.ID
	}
	_ = d.Set("roles", &roleIDs)
	return nil
}

func resourceIAMGroupUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	var group iam.Group
	group.ID = d.Id()
	if d.HasChange("description") {
		group.Description = d.Get("description").(string)
		_, _, err := client.Groups.UpdateGroup(group)
		if err != nil {
			return err
		}
	}
	// Users
	if d.HasChange("users") {
		o, n := d.GetChange("users")
		old := expandStringList(o.(*schema.Set).List())
		newList := expandStringList(n.(*schema.Set).List())
		toAdd := difference(newList, old)
		toRemove := difference(old, newList)

		if len(toRemove) > 0 {
			_, _, _ = client.Groups.RemoveMembers(group, toRemove...)
		}
		if len(toAdd) > 0 {
			_, _, _ = client.Groups.AddMembers(group, toAdd...)
		}
	}

	// Services
	if d.HasChange("services") {
		o, n := d.GetChange("services")
		old := expandStringList(o.(*schema.Set).List())
		newList := expandStringList(n.(*schema.Set).List())
		toAdd := difference(newList, old)
		toRemove := difference(old, newList)

		if len(toRemove) > 0 {
			_, _, _ = client.Groups.RemoveServices(group, toRemove...)
		}
		if len(toAdd) > 0 {
			_, _, _ = client.Groups.AddServices(group, toAdd...)
		}
	}

	if d.HasChange("roles") {
		o, n := d.GetChange("roles")
		old := expandStringList(o.(*schema.Set).List())
		new := expandStringList(n.(*schema.Set).List())
		toAdd := difference(new, old)
		toRemove := difference(old, new)

		// Handle additions
		if len(toAdd) > 0 {
			for _, v := range toAdd {
				var role = iam.Role{ID: v}
				_, _, err := client.Groups.AssignRole(group, role)
				if err != nil {
					return err
				}
			}
		}

		// Remove every role. Simpler to remove and add new ones,
		for _, v := range toRemove {
			var role = iam.Role{ID: v}
			_, _, err := client.Groups.RemoveRole(group, role)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func resourceIAMGroupDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	var group iam.Group
	group.ID = d.Id()

	// Remove all (known) users first before attempting delete
	users := expandStringList(d.Get("users").(*schema.Set).List())
	if len(users) > 0 {
		_, _, err := client.Groups.RemoveMembers(group, users...)
		if err != nil {
			return err
		}
	}

	// Remove all (known) users first before attempting delete
	services := expandStringList(d.Get("services").(*schema.Set).List())
	if len(services) > 0 {
		_, _, err := client.Groups.RemoveServices(group, services...)
		if err != nil {
			return err
		}
	}

	ok, _, err := client.Groups.DeleteGroup(group)
	if !ok {
		return err
	}
	d.SetId("")
	return nil
}
