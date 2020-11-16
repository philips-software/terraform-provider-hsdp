package hsdp

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	username := d.Get("username").(string)

	uuid, _, err := client.Users.GetUserIDByLoginID(username)

	if err != nil {
		// Fallback to legacy user find
		uuid, _, err = client.Users.LegacyGetUserIDByLoginID(username)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(uuid)
	_ = d.Set("uuid", uuid)

	return diags
}
