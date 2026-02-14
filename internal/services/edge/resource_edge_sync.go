package edge

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-dip-api/stl"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceEdgeSync() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Description:   `The ` + "`hsdp_edge_sync`" + ` resource syncs device discovery to the actual device.`,
		CreateContext: resourceEdgeSyncCreate,
		ReadContext:   resourceEdgeSyncRead,
		DeleteContext: resourceEdgeSyncDelete,

		Schema: map[string]*schema.Schema{
			"triggers": {
				Description: "A map of arbitrary strings that, when changed, will force the 'hsdp_edge_sync' resource to be replaced, re-sync conifg with the device.",
				Type:        schema.TypeMap,
				Required:    true,
				ForceNew:    true,
			},
			"principal": config.PrincipalSchema(),
			"serial_number": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceEdgeSyncDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{}
}

func resourceEdgeSyncRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceEdgeSyncCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics
	var client *stl.Client
	var err error

	principal := config.SchemaToPrincipal(d, m)
	client, err = c.STLClient(principal)

	if err != nil {
		return diag.FromErr(err)
	}
	serialNumber := d.Get("serial_number").(string)
	err = client.Devices.SyncDeviceConfig(ctx, serialNumber)
	if err != nil {
		return diag.FromErr(fmt.Errorf("hsdp_edge_sync: %w", err))
	}
	d.SetId(serialNumber)
	return diags
}
