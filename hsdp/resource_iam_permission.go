package hsdp

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
)

func resourceIAMPermission() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceIAMPermissionCreate,
		Read:   resourceIAMPermissionRead,
		Update: resourceIAMPermissionUpdate,
		Delete: resourceIAMPermissionDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"category": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIAMPermissionCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIAMPermissionRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	id := d.Id()
	permission, resp, err := client.Permissions.GetPermissionByName(id) // NOTE: ID = name
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("category", permission.Category)
	d.Set("description", permission.Description)
	d.Set("name", permission.Name)
	d.Set("type", permission.Type)
	d.SetId(permission.Name)
	d.Set("_id", permission.ID)
	return nil
}

func resourceIAMPermissionUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIAMPermissionDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
