package proposition

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-dip-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceIAMProposition() *schema.Resource {
	return &schema.Resource{
		Description: descriptions["proposition"],
		ReadContext: dataSourceIAMPropositionRead,
		Schema: map[string]*schema.Schema{
			"proposition_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organization_id": {
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

func dataSourceIAMPropositionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	propositionID := d.Get("proposition_id").(string)
	name := d.Get("name").(string)
	orgID := d.Get("organization_id").(string)

	var prop *iam.Proposition

	if propositionID != "" {
		prop, _, err = client.Propositions.GetPropositionByID(propositionID)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		if name == "" || orgID == "" {
			return diag.Errorf("when proposition_id is not provided, both name and organization_id are required")
		}
		prop, _, err = client.Propositions.GetProposition(&iam.GetPropositionsOptions{
			OrganizationID: &orgID,
			Name:           &name,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(prop.ID)
	_ = d.Set("name", prop.Name)
	_ = d.Set("description", prop.Description)
	_ = d.Set("global_reference_id", prop.GlobalReferenceID)
	return diags
}
