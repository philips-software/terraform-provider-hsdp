package mdm

import (
	"context"
	"net/http"

	"github.com/philips-software/go-dip-api/connect/mdm"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func DataSourceConnectMDMRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMRegionsRead,
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
			"hsdp_enabled": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
			},
		},
	}

}

func dataSourceConnectMDMRegionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var resources *[]mdm.Region
	var resp *mdm.Response

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		resources, resp, err = client.Regions.GetRegions(nil)
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []string
	var names []string
	var descriptions []string
	var categories []string
	var hsdpEnabled []bool

	for _, r := range *resources {
		ids = append(ids, r.ID)
		names = append(names, r.Name)
		descriptions = append(descriptions, r.Description)
		categories = append(categories, r.Category)
		hsdpEnabled = append(hsdpEnabled, r.HsdpEnabled)
	}
	_ = d.Set("ids", ids)
	_ = d.Set("names", names)
	_ = d.Set("descriptions", descriptions)
	_ = d.Set("categories", categories)
	_ = d.Set("hsdp_enabled", hsdpEnabled)
	result, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(result)
	return diags
}
