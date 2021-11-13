package mdm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceResourceLimits() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMResourcesLimitsRead,
		Schema: map[string]*schema.Schema{
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"defaults": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"overrides": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}

}

func dataSourceConnectMDMResourcesLimitsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	defaults, _, err := client.ResourcesLimits.GetDefault()
	if err != nil {
		return diag.FromErr(err)
	}
	overrides, _, err := client.ResourcesLimits.GetOverride()
	if err != nil {
		return diag.FromErr(err)
	}
	if len(*defaults) == 0 {
		return diags
	}
	var resources []string
	var defaultLimits []int
	var overrideLimits []int
	for k, v := range *defaults {
		resources = append(resources, k)
		defaultLimits = append(defaultLimits, v)
		if l, ok := (*overrides)[k]; ok {
			overrideLimits = append(overrideLimits, l)
		} else {
			overrideLimits = append(overrideLimits, v)
		}
	}
	d.SetId(fmt.Sprintf("DefaultLimits"))
	_ = d.Set("resources", resources)
	_ = d.Set("defaults", defaultLimits)
	_ = d.Set("overrides", overrideLimits)
	return diags
}
