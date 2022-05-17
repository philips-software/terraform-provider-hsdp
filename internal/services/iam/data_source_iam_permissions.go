package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceIAMPermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIAMPermissionsRead,
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
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"categories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"permissions": {
				Type:       schema.TypeList,
				Computed:   true,
				Deprecated: "Use the 'names' field",
				Elem:       &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceIAMPermissionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resp, _, err := client.Permissions.GetPermissions(nil) // Get all permissions

	if err != nil {
		return diag.FromErr(err)
	}
	var ids []string
	var names []string
	var types []string
	var categories []string
	var descriptions []string
	var permissions []string

	for _, p := range *resp {
		ids = append(ids, p.ID)
		names = append(names, p.Name)
		categories = append(categories, p.Category)
		descriptions = append(descriptions, p.Description)
		types = append(types, p.Type)
	}
	d.SetId("permissions")
	_ = d.Set("names", names)
	_ = d.Set("permissions", permissions)
	_ = d.Set("types", types)
	_ = d.Set("categories", categories)
	_ = d.Set("descriptions", descriptions)
	_ = d.Set("ids", ids)

	return diags
}
