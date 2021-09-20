package hsdp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAIWorkspaceInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAIWorkspaceInstanceRead,
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

func dataSourceAIWorkspaceInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var diags diag.Diagnostics

	baseURL := d.Get("base_url").(string)
	InferenceOrgID := d.Get("organization_id").(string)

	client, err := config.getAIWorkspaceClient(baseURL, InferenceOrgID)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	d.SetId(baseURL + InferenceOrgID)
	_ = d.Set("endpoint", client.GetEndpointURL())
	return diags
}
