package hsdp

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIAMIntrospect() *schema.Resource {
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

func dataSourceIAMIntrospectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
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
	_ = d.Set("token", client.Token())
	_ = d.Set("introspect", string(introspectJSON))

	return diags
}
