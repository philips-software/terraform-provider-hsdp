package fhir_store

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceCDRFHIRStore() *schema.Resource {
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

func dataSourceCDRFHIRStoreRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	baseURL := d.Get("base_url").(string)
	fhirOrgID := d.Get("fhir_org_id").(string)
	storeType := "EHR"

	if strings.HasSuffix(baseURL, "/") {
		return diag.FromErr(fmt.Errorf("the base_url should not end with a '/'"))
	}

	if !strings.HasSuffix(baseURL, "/store/fhir") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "the 'base_url' should include the /store/fhir path",
			Detail:   "Please append '/store/fhir' to the base_url to remove this warning",
		})
		baseURL += "/store/fhir"
	}

	client, err := c.GetFHIRClient(baseURL, fhirOrgID)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	d.SetId(fmt.Sprintf("%s:%s/%s", storeType, baseURL, fhirOrgID))
	_ = d.Set("endpoint", client.GetEndpointURL())
	_ = d.Set("type", storeType)
	return diags
}
