package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hsdp/go-hsdp-iam/api"
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
				Computed: true,
			},
			"category": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
	client := m.(*api.Client)

	id := d.Id()
	permission, _, err := client.Permissions.GetPermissionByName(id) // NOTE: ID = name
	if err != nil {
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

func resourceIAMPermissionImport(d *schema.ResourceData, m interface{}) error {
	return nil
}
