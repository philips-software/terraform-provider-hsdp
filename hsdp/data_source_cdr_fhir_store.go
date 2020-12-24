package hsdp

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCDRFHIRStore() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCDRFHIRStoreRead,
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fhir_org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceCDRFHIRStoreRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var diags diag.Diagnostics

	baseURL := d.Get("base_url").(string)
	fhirOrgID := d.Get("fhir_org_id").(string)

	client, err := config.getFHIRClient(baseURL, fhirOrgID)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	d.SetId(baseURL)
	_ = d.Set("endpoint", client.GetEndpointURL())
	_ = d.Set("type", "EHR")
	return diags
}
