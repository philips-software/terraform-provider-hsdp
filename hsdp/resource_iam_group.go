package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
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
		},
	}
}

func resourceIAMGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*iam.Client)
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

	for _, r := range roles {
		role, _, _ := client.Roles.GetRoleByID(r)
		if role != nil {
			client.Groups.AssignRole(*createdGroup, *role)
		}
	}
	return nil
}

func resourceIAMGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*iam.Client)

	id := d.Id()
	group, _, err := client.Groups.GetGroupByID(id)
	if err != nil {
		return err
	}
	d.Set("managing_organization", group.ManagingOrganization)
	d.Set("description", group.Description)
	d.Set("name", group.Name)
	roles, _, err := client.Groups.GetRoles(*group)
	if err != nil {
		return err
	}
	roleIDs := make([]string, len(*roles))
	for i, r := range *roles {
		roleIDs[i] = r.ID
	}
	d.Set("roles", &roleIDs)
	return nil
}

func resourceIAMGroupUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*iam.Client)
	var group iam.Group
	group.ID = d.Id()
	d.Partial(true)
	if d.HasChange("description") {
		group.Description = d.Get("description").(string)
		_, _, err := client.Groups.UpdateGroup(group)
		if err != nil {
			return err
		}
	}
	if d.HasChange("roles") {
		o, n := d.GetChange("roles")
		old := expandStringList(o.(*schema.Set).List())
		new := expandStringList(n.(*schema.Set).List())

		// Remove every role. Simpler to remove and add new ones,
		for _, v := range old {
			var role = iam.Role{ID: v}
			_, _, err := client.Groups.RemoveRole(group, role)
			if err != nil {
				return err
			}
		}
		// Handle additions
		if len(new) > 0 {
			for _, v := range new {
				var role = iam.Role{ID: v}
				_, _, err := client.Groups.AssignRole(group, role)
				if err != nil {
					return err
				}
			}
		}
	}
	d.Partial(false)
	return nil
}

func resourceIAMGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*iam.Client)
	var group iam.Group
	group.ID = d.Id()
	ok, _, err := client.Groups.DeleteGroup(group)
	if !ok {
		return err
	}
	d.SetId("")
	return nil
}
