package ch

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func DataSourceContainerHost() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerHostRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"launch_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"block_devices": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ldap_groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     tools.StringSchema(),
			},
			"protection": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}

}

func dataSourceContainerHostRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	cfg := m.(*config.Config)
	client, err := cfg.CartelClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)

	instance, _, err := client.GetDetails(name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instance.InstanceID)

	_ = d.Set("type", instance.InstanceType)
	_ = d.Set("owner", instance.Owner)
	_ = d.Set("private_ip", instance.PrivateAddress)
	_ = d.Set("public_ip", instance.PublicAddress)
	_ = d.Set("state", instance.State)
	_ = d.Set("state", instance.State)
	_ = d.Set("launch_time", instance.LaunchTime)
	_ = d.Set("block_devices", tools.SchemaSetStrings(instance.BlockDevices))
	_ = d.Set("security_groups", tools.SchemaSetStrings(instance.SecurityGroups))
	_ = d.Set("ldap_groups", tools.SchemaSetStrings(instance.LdapGroups))
	_ = d.Set("role", instance.Role)
	_ = d.Set("subnet", instance.Subnet)
	_ = d.Set("vpc", instance.Vpc)
	_ = d.Set("zone", instance.Zone)
	_ = d.Set("tags", instance.Tags)
	_ = d.Set("protection", instance.Protection)
	return diags
}
