package inference

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/ai"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func DataSourceAIInferenceComputeTargets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAIInferenceComputeTargetsRead,
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

func dataSourceAIInferenceComputeTargetsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	endpoint := d.Get("endpoint").(string)
	client, err := c.GetAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	var targets []ai.ComputeTarget
	err = tools.TryAICall(func() (*ai.Response, error) {
		var err error
		var resp *ai.Response
		_ = client.TokenRefresh()
		targets, resp, err = client.ComputeTarget.GetComputeTargets(nil)
		return resp, err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("ai_inference_compute_targets")

	var names []string
	var ids []string

	for _, env := range targets {
		names = append(names, env.Name)
		ids = append(ids, env.ID)
	}
	_ = d.Set("names", names)
	_ = d.Set("ids", ids)

	return diags
}
