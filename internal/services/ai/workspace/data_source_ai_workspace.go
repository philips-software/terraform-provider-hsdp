package workspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceAIWorkspace() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source is deprecated and will be removed in an upcoming release of the provider",
		ReadContext:        dataSourceAIWorkspaceRead,
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
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	baseURL := d.Get("base_url").(string)
	workspaceTenantID := d.Get("workspace_tenant_id").(string)
	workspaceID := d.Get("workspace_id").(string)

	client, err := c.GetAIWorkspaceClient(baseURL, workspaceTenantID)
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
