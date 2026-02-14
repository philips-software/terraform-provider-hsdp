package mdm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-dip-api/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceConnectMDMRegion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMRegionRead,
		Schema: map[string]*schema.Schema{
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hsdp_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}

}

func dataSourceConnectMDMRegionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)

	resources, _, err := client.Regions.GetRegions(&mdm.GetRegionOptions{
		Name: &name,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if len(*resources) == 0 {
		return diag.FromErr(fmt.Errorf("region '%s' not found", name))
	}
	resource := (*resources)[0]
	d.SetId(fmt.Sprintf("Region/%s", resource.ID))
	_ = d.Set("guid", resource.ID)
	_ = d.Set("name", resource.Name)
	_ = d.Set("description", resource.Description)
	_ = d.Set("category", resource.Category)
	_ = d.Set("hsdp_enabled", resource.HsdpEnabled)
	return diags
}
