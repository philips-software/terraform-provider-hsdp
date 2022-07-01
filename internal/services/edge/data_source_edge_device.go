package edge

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/stl"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceEdgeDevice() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEdgeDeviceRead,
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"serial_number": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hardware_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_interface_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceEdgeDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics
	var client *stl.Client
	var err error

	endpoint := d.Get("endpoint").(string)
	if endpoint != "" {
		client, err = c.STLClient(&config.Principal{Endpoint: endpoint})
	} else {
		client, err = c.STLClient()
	}
	if err != nil {
		return diag.FromErr(err)
	}
	serialNumber := d.Get("serial_number").(string)
	device, err := client.Devices.GetDeviceBySerial(ctx, serialNumber)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read Edge device: %w", err))
	}
	_ = d.Set("region", device.Region)
	_ = d.Set("name", device.Name)
	_ = d.Set("state", device.State)
	_ = d.Set("primary_interface_ip", device.PrimaryInterface.Address)
	d.SetId(fmt.Sprintf("%d", device.ID))
	return diags
}
