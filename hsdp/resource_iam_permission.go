package hsdp

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	id := d.Id()
	permission, resp, err := client.Permissions.GetPermissionByName(id) // NOTE: ID = name
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	_ = d.Set("category", permission.Category)
	_ = d.Set("description", permission.Description)
	_ = d.Set("name", permission.Name)
	_ = d.Set("type", permission.Type)
	_ = d.Set("_id", permission.ID)
	d.SetId(permission.Name)
	return nil
}

func resourceIAMPermissionUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIAMPermissionDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
