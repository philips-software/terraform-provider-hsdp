package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/loafoe/go-hsdp/api"
)

func resourceIAMRole() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceIAMRoleCreate,
		Read:   resourceIAMRoleRead,
		Update: resourceIAMRoleUpdate,
		Delete: resourceIAMRoleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"managing_organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIAMRoleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*api.Client)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	managingOrganization := d.Get("managing_organization").(string)
	permissions := expandStringList(d.Get("permissions").(*schema.Set).List())

	role, _, err := client.Roles.CreateRole(name, description, managingOrganization)
	if err != nil {
		return err
	}
	for _, p := range permissions {
		client.Roles.AddRolePermission(*role, p)
	}
	d.SetId(role.ID)
	return nil
}

func resourceIAMRoleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*api.Client)

	id := d.Id()
	role, _, err := client.Roles.GetRoleByID(id)
	if err != nil {
		return err
	}
	d.Set("description", role.Description)
	d.Set("name", role.Name)
	d.Set("managing_organization", role.ManagingOrganization)
	d.SetId(role.ID)

	permissions, err := client.Roles.GetRolePermissions(*role)
	if err != nil {
		return err
	}
	d.Set("permissions", permissions)
	return nil
}

func resourceIAMRoleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*api.Client)
	id := d.Id()
	role, _, err := client.Roles.GetRoleByID(id)
	if err != nil {
		return err
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		role.Description = description
		client.Roles.UpdateRole(role)
	}
	if d.HasChange("permissions") {
		o, n := d.GetChange("permissions")
		old := expandStringList(o.(*schema.Set).List())
		new := expandStringList(n.(*schema.Set).List())

		// Remove every permission. Simpler to remove and add new ones,
		for _, v := range old {
			_, _, err := client.Roles.RemoveRolePermission(*role, v)
			if err != nil {
				return err
			}
		}

		// Handle additions
		if len(new) > 0 {
			for _, v := range new {
				_, _, err := client.Roles.AddRolePermission(*role, v)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func resourceIAMRoleDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIAMRoleImport(d *schema.ResourceData, m interface{}) error {
	return nil
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}
