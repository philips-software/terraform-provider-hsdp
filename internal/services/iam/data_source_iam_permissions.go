package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	config2 "github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceIAMPermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIAMPermissionsRead,
		Schema: map[string]*schema.Schema{
			"permissions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}

}

func dataSourceIAMPermissionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config2.Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resp, _, err := client.Permissions.GetPermissions(nil) // Get all permissions

	if err != nil {
		return diag.FromErr(err)
	}
	permissions := make([]string, 0)

	for _, p := range *resp {
		permissions = append(permissions, p.Name)
	}
	d.SetId("permissions")
	_ = d.Set("permissions", permissions)

	return diags
}
