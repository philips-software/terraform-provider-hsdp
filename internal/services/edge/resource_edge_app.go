package edge

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/stl"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceEdgeApp() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceEdgeAppCreate,
		ReadContext:   resourceEdgeAppRead,
		UpdateContext: resourceEdgeAppUpdate,
		DeleteContext: resourceEdgeAppDelete,

		Schema: map[string]*schema.Schema{
			"serial_number": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"device_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sync": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceEdgeAppUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics
	var client *stl.Client
	var err error

	if endpoint, ok := d.GetOk("endpoint"); ok {
		client, err = c.STLClient(endpoint.(string))
	} else {
		client, err = c.STLClient()
	}
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)
	content := d.Get("content").(string)
	var resourceID int64
	_, _ = fmt.Sscanf(d.Id(), "%d", &resourceID)
	_, err = client.Apps.UpdateAppResource(ctx, stl.UpdateApplicationResourceInput{
		ID:       resourceID,
		Name:     name,
		Content:  base64.StdEncoding.EncodeToString([]byte(content)),
		IsLocked: false,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("edge_app: update Edge app: %w", err))
	}
	syncSTLIfNeeded(ctx, client, d, m)
	return diags
}

func resourceEdgeAppDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics
	var client *stl.Client
	var err error

	if endpoint, ok := d.GetOk("endpoint"); ok {
		client, err = c.STLClient(endpoint.(string))
	} else {
		client, err = c.STLClient()
	}
	if err != nil {
		return diag.FromErr(err)
	}
	var resourceID int64
	_, _ = fmt.Sscanf(d.Id(), "%d", &resourceID)
	_, err = client.Apps.DeleteAppResource(ctx, stl.DeleteApplicationResourceInput{
		ID: resourceID,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("edge_app: delete Edge resource: %w", err))
	}
	syncSTLIfNeeded(ctx, client, d, m)
	d.SetId("")
	return diags
}

func resourceEdgeAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics
	var client *stl.Client
	var err error

	if endpoint, ok := d.GetOk("endpoint"); ok {
		client, err = c.STLClient(endpoint.(string))
	} else {
		client, err = c.STLClient()
	}
	if err != nil {
		return diag.FromErr(err)
	}
	var resourceID int64
	_, _ = fmt.Sscanf(d.Id(), "%d", &resourceID)
	resource, err := client.Apps.GetAppResourceByID(ctx, resourceID)
	if err != nil {
		if strings.Contains(err.Error(), "resource not found") { // GraphQL downside
			d.SetId("")
			return diags
		}
		return diag.FromErr(fmt.Errorf("edge_app: read Edge device: %w", err))
	}
	_ = d.Set("name", resource.Name)
	content, err := base64.StdEncoding.DecodeString(resource.Content)
	if err != nil {
		return diag.FromErr(fmt.Errorf("edge_app: decode content: %w", err))
	}
	_ = d.Set("content", string(content))
	_ = d.Set("device_id", resource.DeviceID)
	device, err := client.Devices.GetDeviceByID(ctx, resource.DeviceID)
	if err == nil {
		_ = d.Set("serial_number", device.SerialNumber)
	}
	return diags
}

func resourceEdgeAppCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var client *stl.Client
	var err error

	if endpoint, ok := d.GetOk("endpoint"); ok {
		client, err = c.STLClient(endpoint.(string))
	} else {
		client, err = c.STLClient()
	}
	if err != nil {
		return diag.FromErr(err)
	}
	serialNumber := d.Get("serial_number").(string)
	name := d.Get("name").(string)
	content := d.Get("content").(string)
	resource, err := client.Apps.CreateAppResource(ctx, stl.CreateApplicationResourceInput{
		SerialNumber: serialNumber,
		Name:         name,
		Content:      base64.StdEncoding.EncodeToString([]byte(content)),
		IsLocked:     false,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("edge_app: create Edge app: %w", err))
	}
	d.SetId(fmt.Sprintf("%d", resource.ID))
	syncSTLIfNeeded(ctx, client, d, m)
	return resourceEdgeAppRead(ctx, d, m)
}
