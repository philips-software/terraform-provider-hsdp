package hsdp

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func dataSourceIAMApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIAMApplicationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"proposition_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"global_reference_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceIAMApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	propID := d.Get("proposition_id").(string)
	name := d.Get("name").(string)

	prop, _, err := client.Applications.GetApplication(&iam.GetApplicationsOptions{
		PropositionID: &propID,
		Name:          &name,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if len(prop) == 0 {
		return diag.FromErr(ErrResourceNotFound)
	}

	d.SetId(prop[0].ID)
	_ = d.Set("description", prop[0].Description)
	_ = d.Set("global_reference_id", prop[0].GlobalReferenceID)
	return diags
}
