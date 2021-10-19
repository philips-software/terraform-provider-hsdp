package inference

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceAIInferenceComputeEnvironments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAIInferenceComputeEnvironmentsRead,
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceAIInferenceComputeEnvironmentsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	endpoint := d.Get("endpoint").(string)
	client, err := c.GetAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	environments, _, err := client.ComputeEnvironment.GetComputeEnvironments(nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("ai_inference_compute_environments")

	names := make([]string, len(environments))
	ids := make([]string, len(environments))

	for i := 0; i < len(environments); i++ {
		names[i] = environments[i].Name
		ids[i] = environments[i].ID
	}
	_ = d.Set("names", names)
	_ = d.Set("ids", ids)

	return diags
}
