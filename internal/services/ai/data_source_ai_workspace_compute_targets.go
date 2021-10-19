package ai

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceAIWorkspaceComputeTargets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAIWorkspaceComputeTargetsRead,
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

func dataSourceAIWorkspaceComputeTargetsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	endpoint := d.Get("endpoint").(string)
	client, err := c.GetAIWorkspaceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	environments, _, err := client.ComputeTarget.GetComputeTargets(nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("ai_workspace_compute_targets")

	var names []string
	var ids []string

	for _, env := range environments {
		names = append(names, env.Name)
		ids = append(ids, env.ID)
	}
	_ = d.Set("names", names)
	_ = d.Set("ids", ids)

	return diags
}
