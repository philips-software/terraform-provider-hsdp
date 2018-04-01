package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
	iamclient "github.com/loafoe/go-hsdpiam"
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
			"org_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceIAMGroupCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIAMGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*iamclient.Client)

	id := d.Id()
	group, _, err := client.Groups.GetGroup(&iamclient.GetGroupOptions{ID: &id})
	if err != nil {
		return err
	}
	d.Set("org_id", group.OrganizationID)
	d.Set("description", group.Description)
	d.Set("name", group.Name)
	return nil
}

func resourceIAMGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIAMGroupDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIAMGroupImport(d *schema.ResourceData, m interface{}) error {
	return nil
}
