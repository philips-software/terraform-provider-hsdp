package hsdp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/cdl"
)

func dataSourceCDLResearchStudies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCDLResearchStudiesRead,
		Schema: map[string]*schema.Schema{
			"cdl_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"titles": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceCDLResearchStudiesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	endpoint := d.Get("cdl_endpoint").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
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
