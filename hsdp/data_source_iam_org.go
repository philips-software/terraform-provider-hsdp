package hsdp

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIAMOrg() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIAMOrgRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceIAMOrgRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}
	orgId := d.Get("organization_id").(string)

	org, _, err := client.Organizations.GetOrganizationByID(orgId) // Get all permissions

	if err != nil {
		return err
	}

	d.SetId(orgId)
	_ = d.Set("name", org.Name)
	_ = d.Set("description", org.Description)
	_ = d.Set("active", org.Active)
	_ = d.Set("type", org.Type)
	_ = d.Set("external_id", org.ExternalID)
	_ = d.Set("display_name", org.DisplayName)
	_ = d.Set("parent_org_id", org.Parent.Value)

	return err
}
