package hsdp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCDLStudy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCDLStudyRead,
		Schema: map[string]*schema.Schema{
			"study_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cdl_instance_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ends_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"study_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceCDLStudyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
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
