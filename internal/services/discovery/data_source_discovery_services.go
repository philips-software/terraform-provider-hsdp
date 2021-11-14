package discovery

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceDiscoveryServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDiscoveryServicesRead,
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
			"trusted": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
			},
		},
	}

}

func dataSourceDiscoveryServicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.DiscoveryClient()
	if err != nil {
		return diag.FromErr(err)
	}

	services, _, err := client.GetServices()
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []string
	var names []string
	var tags []string
	var trusted []bool

	for _, s := range *services {
		// All criteria match, so add user
		ids = append(ids, s.ID)
		names = append(names, s.Name)
		tags = append(tags, s.Tag)
		trusted = append(trusted, s.IsTrusted)
	}
	_ = d.Set("ids", ids)
	_ = d.Set("names", names)
	_ = d.Set("tags", tags)
	_ = d.Set("trusted", trusted)
	result, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(result)
	return diags
}
