package ch

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceContainerHostInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceContainerHostInstancesRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"owners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceContainerHostInstancesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	cfg := m.(*config.Config)
	client, err := cfg.CartelClient()
	if err != nil {
		return diag.FromErr(err)
	}

	instances, _, err := client.GetAllInstances()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("cartel_instances")

	var names []string
	var ids []string
	var roles []string
	var owners []string

	for _, instance := range *instances {
		names = append(names, instance.NameTag)
		ids = append(ids, instance.InstanceID)
		roles = append(roles, instance.Role)
		owners = append(owners, instance.Owner)
	}
	_ = d.Set("names", names)
	_ = d.Set("ids", ids)
	_ = d.Set("owners", owners)
	_ = d.Set("roles", roles)

	return diags
}
