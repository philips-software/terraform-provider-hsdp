package ch

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceContainerHostSecurityGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerHostSecurityGroupsRead,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}

}

func dataSourceContainerHostSecurityGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	var diags diag.Diagnostics

	client, err := c.CartelClient()
	if err != nil {
		return diag.FromErr(err)
	}
	details, _, err := client.GetSecurityGroups()
	if err != nil {
		return diag.FromErr(err)
	}

	names := make([]string, 0)
	for _, name := range *details {
		names = append(names, name)
	}
	_ = d.Set("names", names)
	d.SetId("names")
	return diags
}
