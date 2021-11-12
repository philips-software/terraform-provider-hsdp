package mdm

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceConnectMDMOauthClientScopes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMOAuthClientScopesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"organizations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"propositions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bootstrap_enabled": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
			},
		},
	}

}

func dataSourceConnectMDMOAuthClientScopesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resourcesList, _, err := client.OAuthClientScopes.GetOAuthClientScopes(nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []string
	var organizations []string
	var propositions []string
	var actions []string
	var services []string
	var resources []string
	var bootstrapEnabled []bool

	for _, r := range *resourcesList {
		ids = append(ids, r.ID)
		resources = append(resources, r.Resource)
		propositions = append(propositions, r.Proposition)
		services = append(services, r.Service)
		actions = append(actions, r.Action)
		organizations = append(organizations, r.Organization)
		bootstrapEnabled = append(bootstrapEnabled, r.BootstrapEnabled)
	}
	_ = d.Set("ids", ids)
	_ = d.Set("organizations", organizations)
	_ = d.Set("propositions", propositions)
	_ = d.Set("services", services)
	_ = d.Set("actions", actions)
	_ = d.Set("resources", resources)
	_ = d.Set("bootstrap_enabled", bootstrapEnabled)

	result, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(result)
	return diags
}
