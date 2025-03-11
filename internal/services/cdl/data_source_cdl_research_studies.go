package cdl

import (
	"context"

	"github.com/dip-software/go-dip-api/cdl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceCDLResearchStudies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCDLResearchStudiesRead,
		Schema: map[string]*schema.Schema{
			"cdl_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"titles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceCDLResearchStudiesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	endpoint := d.Get("cdl_endpoint").(string)

	client, err := c.GetCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	studies, _, err := client.Study.GetStudies(&cdl.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(endpoint + "studies")

	var titles []string
	var ids []string

	for _, study := range studies {
		titles = append(titles, study.Title)
		ids = append(ids, study.ID)
	}
	_ = d.Set("titles", titles)
	_ = d.Set("ids", ids)

	return diags
}
