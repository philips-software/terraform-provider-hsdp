package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceIAMToken() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIAMTokenRead,
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"id_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceIAMTokenRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("token-" + client.BaseIAMURL().Host)
	token, err := client.Token()
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("access_token", token)
	_ = d.Set("expires_at", client.Expires())
	_ = d.Set("id_token", client.IDToken())

	return diags
}
