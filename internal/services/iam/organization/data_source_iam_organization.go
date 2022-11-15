package organization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceIAMOrg() *schema.Resource {
	return &schema.Resource{
		Description: descriptions["organization"],
		ReadContext: dataSourceIAMOrgRead,
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

func dataSourceIAMOrgRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	orgId := d.Get("organization_id").(string)

	org, _, err := client.Organizations.GetOrganizationByID(orgId) // Get all permissions

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(orgId)
	_ = d.Set("name", org.Name)
	_ = d.Set("description", org.Description)
	_ = d.Set("active", org.Active)
	_ = d.Set("type", org.Type)
	_ = d.Set("external_id", org.ExternalID)
	_ = d.Set("display_name", org.DisplayName)
	_ = d.Set("parent_org_id", org.Parent.Value)

	return diags
}
