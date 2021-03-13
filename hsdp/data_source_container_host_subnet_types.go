package hsdp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContainerHostSubnetTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerHostSubnetTypesRead,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ids": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"networks": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}

}

func dataSourceContainerHostSubnetTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	var diags diag.Diagnostics

	client, err := config.CartelClient()
	if err != nil {
		return diag.FromErr(err)
	}
	details, _, err := client.GetAllSubnets()
	if err != nil {
		return diag.FromErr(err)
	}

	ids := make(map[string]interface{})
	networks := make(map[string]interface{})
	names := make([]string, 0)
	for name, subnet := range *details {
		names = append(names, name)
		ids[name] = subnet.ID
		networks[name] = subnet.Network
	}
	_ = d.Set("ids", ids)
	_ = d.Set("networks", networks)
	_ = d.Set("names", names)
	d.SetId("subnets")
	return diags
}
