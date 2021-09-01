package hsdp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContainerHostInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerHostInstancesRead,
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
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"owners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"private_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceContainerHostInstancesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	client, err := config.CartelClient()
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
	var types []string
	var privateIPs []string
	var owners []string

	for _, instance := range *instances {
		names = append(names, instance.NameTag)
		ids = append(ids, instance.InstanceID)
		types = append(types, instance.InstanceType)
		privateIPs = append(privateIPs, instance.PrivateAddress)
		owners = append(owners, instance.Owner)
	}
	_ = d.Set("names", names)
	_ = d.Set("ids", ids)
	_ = d.Set("types", types)
	_ = d.Set("owners", owners)
	_ = d.Set("private_ips", privateIPs)

	return diags
}
