package application

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-dip-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceIAMApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIAMApplicationRead,
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proposition_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global_reference_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceIAMApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	applicationID := d.Get("application_id").(string)
	name := d.Get("name").(string)
	propID := d.Get("proposition_id").(string)

	var app *iam.Application

	if applicationID != "" {
		app, _, err = client.Applications.GetApplicationByID(applicationID)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		if name == "" || propID == "" {
			return diag.Errorf("when application_id is not provided, both name and proposition_id are required")
		}
		apps, _, err := client.Applications.GetApplications(&iam.GetApplicationsOptions{
			PropositionID: &propID,
			Name:          &name,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		if len(apps) == 0 {
			return diag.FromErr(config.ErrResourceNotFound)
		}
		app = apps[0]
	}

	d.SetId(app.ID)
	_ = d.Set("name", app.Name)
	_ = d.Set("description", app.Description)
	_ = d.Set("global_reference_id", app.GlobalReferenceID)
	return diags
}
