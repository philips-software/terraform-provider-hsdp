package mdm

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-dip-api/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceConnectMDMProposition() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMPropositionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"error_on_not_found": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"found": {
				Type:     schema.TypeBool,
				Computed: true,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proposition_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proposition_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceConnectMDMPropositionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	orgId := d.Get("organization_id").(string)
	name := d.Get("name").(string)
	errorOnNotFound := d.Get("error_on_not_found").(bool)

	prop, resp, err := client.Propositions.GetProposition(&mdm.GetPropositionsOptions{
		OrganizationID: &orgId,
		Name:           &name,
	})
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		if (errors.Is(err, mdm.ErrEmptyResult) || resp.StatusCode() == http.StatusNotFound) && errorOnNotFound {
			return diag.FromErr(err)
		}
		// Not found, but no error
		_ = d.Set("found", false)
		d.SetId("PropositionNotFound")
		return diags
	}

	d.SetId(fmt.Sprintf("Proposition/%s", prop.ID))
	propositionToSchema(*prop, d)
	_ = d.Set("found", true)
	return diags
}
