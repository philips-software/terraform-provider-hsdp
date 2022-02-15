package iam

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIAMIntrospect() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIAMIntrospectRead,
		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
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

	resp, _, err := client.Introspect()

	if err != nil {
		return diag.FromErr(err)
	}
	introspectJSON, err := json.Marshal(&resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Username)
	_ = d.Set("managing_organization", resp.Organizations.ManagingOrganization)
	_ = d.Set("username", resp.Username)
	token, err := client.Token()
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("token", token)
	_ = d.Set("introspect", string(introspectJSON))

	return diags
}
