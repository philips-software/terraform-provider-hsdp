package hsdp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAIWorkspace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAIWorkspaceRead,
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceAIWorkspaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var diags diag.Diagnostics

	baseURL := d.Get("base_url").(string)
	workspaceTenantID := d.Get("workspace_tenant_id").(string)
	workspaceID := d.Get("workspace_id").(string)

	client, err := config.getAIWorkspaceClient(baseURL, workspaceTenantID)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	ws, _, err := client.Workspace.GetWorkspaceByID(workspaceID)
	if err != nil {
		return diag.FromErr(err)
	}
	accessURL, _, err := client.Workspace.GetWorkspaceAccessURL(*ws)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(workspaceID)
	_ = d.Set("endpoint", client.GetEndpointURL())
	_ = d.Set("url", accessURL.URL)
	return diags
}
