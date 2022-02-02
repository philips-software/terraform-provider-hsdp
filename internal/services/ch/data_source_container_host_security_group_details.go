package ch

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceContainerHostSecurityGroupDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerHostSecurityGroupDetailsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocols": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"port_ranges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}

}

func dataSourceContainerHostSecurityGroupDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	var diags diag.Diagnostics

	client, err := c.CartelClient()
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)

	details, _, err := client.GetSecurityGroupDetails(name)
	if err != nil {
		return diag.FromErr(err)
	}
	protocols := make([]string, 0)
	sources := make([]string, 0)
	portRanges := make([]string, 0)

	for _, rule := range *details {
		protocols = append(protocols, rule.Protocol)
		sources = append(sources, strings.Join(rule.Source, ","))
		portRanges = append(sources, rule.PortRange)
	}
	_ = d.Set("protocols", protocols)
	_ = d.Set("port_ranges", portRanges)
	_ = d.Set("sources", sources)
	d.SetId(fmt.Sprintf("details-%s", name))
	return diags
}
