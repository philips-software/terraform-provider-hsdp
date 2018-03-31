package hsdpiam

import (
	"github.com/hashicorp/terraform/helper/schema"
	iamclient "github.com/loafoe/go-hsdpiam"
)

func resourceOrg() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceOrgCreate,
		Read:   resourceOrgRead,
		Update: resourceOrgUpdate,
		Delete: resourceOrgDelete,

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

func resourceOrgCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceOrgRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*iamclient.Client)

	id := d.Get("org_id").(string)
	org, _, err := client.Organizations.GetOrganization(&iamclient.GetOrganizationOptions{ID: &id})
	if err != nil {
		return err
	}
	d.Set("org_id", org.OrganizationID)
	d.Set("description", org.Description)
	d.Set("name", org.Name)
	return nil
}

func resourceOrgUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceOrgDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceOrgImport(d *schema.ResourceData, m interface{}) error {
	return nil
}
