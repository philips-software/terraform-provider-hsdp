package inference

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceAIInferenceServiceInstance() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source is deprecated and will be removed in an upcoming release of the provider",
		ReadContext:        dataSourceAIInferenceServiceInstanceRead,
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceAIInferenceServiceInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	baseURL := d.Get("base_url").(string)
	InferenceOrgID := d.Get("organization_id").(string)

	client, err := c.GetAIInferenceClient(baseURL, InferenceOrgID)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	d.SetId(baseURL + InferenceOrgID)
	_ = d.Set("endpoint", client.GetEndpointURL())
	return diags
}
