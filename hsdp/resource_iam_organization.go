package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/loafoe/go-hsdp/api"
)

func resourceIAMOrg() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceIAMOrgCreate,
		Read:   resourceIAMOrgRead,
		Update: resourceIAMOrgUpdate,
		Delete: resourceIAMOrgDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"distinct_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIAMOrgCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIAMOrgRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*api.Client)

	id := d.Id()
	org, _, err := client.Organizations.GetOrganizationByID(id)
	if err != nil {
		return err
	}
	d.Set("org_id", org.OrganizationID)
	d.Set("description", org.Description)
	d.Set("name", org.Name)
	return nil
}

func resourceIAMOrgUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIAMOrgDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIAMOrgImport(d *schema.ResourceData, m interface{}) error {
	return nil
}
