package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceIAMProposition() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIAMPropositionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
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
		},
	}

}

func dataSourceIAMPropositionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	orgId := d.Get("organization_id").(string)
	name := d.Get("name").(string)

	prop, _, err := client.Propositions.GetProposition(&iam.GetPropositionsOptions{
		OrganizationID: &orgId,
		Name:           &name,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(prop.ID)
	_ = d.Set("description", prop.Description)
	_ = d.Set("global_reference_id", prop.GlobalReferenceID)
	return diags
}
