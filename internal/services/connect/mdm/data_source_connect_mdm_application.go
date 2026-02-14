package mdm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-dip-api/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceConnectMDMApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMApplicationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"proposition_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global_reference_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_guid_system": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_guid_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_group_guid_system": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_group_guid_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceConnectMDMApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	propID := d.Get("proposition_id").(string)
	name := d.Get("name").(string)

	apps, _, err := client.Applications.GetApplications(&mdm.GetApplicationsOptions{
		PropositionID: &propID,
		Name:          &name,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if apps == nil || len(*apps) == 0 {
		return diag.FromErr(config.ErrResourceNotFound)
	}
	app := (*apps)[0]

	d.SetId(fmt.Sprintf("Application/%s", app.ID))
	_ = d.Set("guid", app.ID)
	if app.ApplicationGuid != nil {
		_ = d.Set("application_guid_system", app.ApplicationGuid.System)
		_ = d.Set("application_guid_value", app.ApplicationGuid.Value)
	}
	if app.DefaultGroupGuid != nil {
		_ = d.Set("default_group_guid_system", app.DefaultGroupGuid.System)
		_ = d.Set("default_group_guid_value", app.DefaultGroupGuid.Value)
	}
	_ = d.Set("description", app.Description)
	_ = d.Set("global_reference_id", app.GlobalReferenceID)
	return diags
}
