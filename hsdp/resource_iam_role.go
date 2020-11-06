package hsdp

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateUpperString,
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
			"ticket_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIAMRoleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	managingOrganization := d.Get("managing_organization").(string)
	permissions := expandStringList(d.Get("permissions").(*schema.Set).List())

	role, _, err := client.Roles.CreateRole(name, description, managingOrganization)
	if err != nil {
		return err
	}
	for _, p := range permissions {
		_, _, _ = client.Roles.AddRolePermission(*role, p)
	}
	d.SetId(role.ID)
	return resourceIAMRoleRead(d, meta)
}

func resourceIAMRoleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	id := d.Id()
	role, resp, err := client.Roles.GetRoleByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	d.SetId(role.ID)
	_ = d.Set("description", role.Description)
	_ = d.Set("name", role.Name)
	_ = d.Set("managing_organization", role.ManagingOrganization)

	permissions, _, err := client.Roles.GetRolePermissions(*role)
	if err != nil {
		return err
	}
	_ = d.Set("permissions", permissions)
	return nil
}

func resourceIAMRoleUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	id := d.Id()
	role, _, err := client.Roles.GetRoleByID(id)
	if err != nil {
		return err
	}

	if d.HasChange("description") {
		return fmt.Errorf("description changes are not supported")
	}

	if d.HasChange("permissions") {
		o, n := d.GetChange("permissions")
		oldList := expandStringList(o.(*schema.Set).List())
		newList := expandStringList(n.(*schema.Set).List())
		toAdd := difference(newList, oldList)
		toRemove := difference(oldList, newList)

		// Additions
		if len(toAdd) > 0 {
			for _, v := range toAdd {
				_, _, err := client.Roles.AddRolePermission(*role, v)
				if err != nil {
					return err
				}
			}
		}

		// Removals
		for _, v := range toRemove {
			ticketProtection := d.Get("ticket_protection").(bool)
			if ticketProtection && v == "CLIENT.SCOPES" {
				return fmt.Errorf("Refusing to remove CLIENT.SCOPES permission, set ticket_protection to `false` to override")
			}
			_, _, err := client.Roles.RemoveRolePermission(*role, v)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func resourceIAMRoleDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	var role iam.Role
	role.ID = d.Id()

	ok, _, err := client.Roles.DeleteRole(role)
	if !ok {
		return err
	}
	d.SetId("")
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
