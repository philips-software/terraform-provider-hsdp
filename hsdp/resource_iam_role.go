package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
	iamclient "github.com/loafoe/go-hsdpiam"
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"managing_organization": &schema.Schema{
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

func resourceIAMRoleCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIAMRoleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*iamclient.Client)

	id := d.Id()
	role, _, err := client.Roles.GetRoleByID(id)
	if err != nil {
		return err
	}
	d.Set("description", role.Description)
	d.Set("name", role.Name)
	d.Set("managing_organization", role.ManagingOrganization)
	d.SetId(role.ID)
	return nil
}

func resourceIAMRoleUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIAMRoleDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIAMRoleImport(d *schema.ResourceData, m interface{}) error {
	return nil
}
