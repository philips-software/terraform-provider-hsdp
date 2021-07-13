package hsdp

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCDLInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCDLInstanceRead,
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceCDLInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var diags diag.Diagnostics

	baseURL := d.Get("base_url").(string)
	cdlOrgID := d.Get("organization_id").(string)

	if strings.Contains(baseURL, "/store/cdl/") {
		return diag.FromErr(fmt.Errorf("the base_url argument should not have `/store/cdl/` at the end"))
	}

	client, err := config.getCDLClient(baseURL, cdlOrgID)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	d.SetId(baseURL + cdlOrgID)
	_ = d.Set("endpoint", client.GetEndpointURL())
	return diags
}
