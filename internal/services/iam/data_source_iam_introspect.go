package iam

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func DataSourceIAMIntrospect() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIAMIntrospectRead,
		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organization_context": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"token_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"managing_organization": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"introspect": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"effective_permissions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceIAMIntrospectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	orgContext := d.Get("organization_context").(string)

	var resp *iam.IntrospectResponse
	if orgContext != "" {
		resp, _, err = client.Introspect(iam.WithOrgContext(orgContext))
	} else {
		resp, _, err = client.Introspect()
	}
	if err != nil {
		return diag.FromErr(err)
	}
	introspectJSON, err := json.Marshal(&resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Sub)
	_ = d.Set("managing_organization", resp.Organizations.ManagingOrganization)
	_ = d.Set("username", resp.Username)
	token, err := client.Token()
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("token", token)
	_ = d.Set("token_type", resp.TokenType)
	_ = d.Set("identity_type", resp.IdentityType)
	_ = d.Set("subject", resp.Sub)
	_ = d.Set("issuer", resp.ISS)
	_ = d.Set("introspect", string(introspectJSON))
	if orgContext != "" {
		for _, org := range resp.Organizations.OrganizationList {
			if org.OrganizationID != orgContext {
				continue
			}
			_ = d.Set("effective_permissions", tools.SchemaSetStrings(org.EffectivePermissions))
			break
		}
	}
	return diags
}
