package hsdp

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/stl"
)

func resourceSTLApp() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceSTLAppCreate,
		ReadContext:   resourceSTLAppRead,
		UpdateContext: resourceSTLAppUpdate,
		DeleteContext: resourceSTLAppDelete,

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

func resourceSTLAppUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	var diags diag.Diagnostics
	var client *stl.Client
	var err error

	if endpoint, ok := d.GetOk("endpoint"); ok {
		client, err = config.STLClient(endpoint.(string))
	} else {
		client, err = config.STLClient()
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
		return diag.FromErr(fmt.Errorf("stl_app: update STL app: %w", err))
	}
	syncSTLIfNeeded(ctx, client, d, m)
	return diags
}

func resourceSTLAppDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	var diags diag.Diagnostics
	var client *stl.Client
	var err error

	if endpoint, ok := d.GetOk("endpoint"); ok {
		client, err = config.STLClient(endpoint.(string))
	} else {
		client, err = config.STLClient()
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
		return diag.FromErr(fmt.Errorf("stl_app: delete STL resource: %w", err))
	}
	syncSTLIfNeeded(ctx, client, d, m)
	d.SetId("")
	return diags
}

func resourceSTLAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	var diags diag.Diagnostics
	var client *stl.Client
	var err error

	if endpoint, ok := d.GetOk("endpoint"); ok {
		client, err = config.STLClient(endpoint.(string))
	} else {
		client, err = config.STLClient()
	}
	if err != nil {
		return diag.FromErr(err)
	}
	var resourceID int64
	_, _ = fmt.Sscanf(d.Id(), "%d", &resourceID)
	resource, err := client.Apps.GetAppResourceByID(ctx, resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("stl_app: read STL device: %w", err))
	}
	_ = d.Set("name", resource.Name)
	content, err := base64.StdEncoding.DecodeString(resource.Content)
	if err != nil {
		return diag.FromErr(fmt.Errorf("stl_app: decode content: %w", err))
	}
	_ = d.Set("content", content)
	_ = d.Set("device_id", resource.DeviceID)
	device, err := client.Devices.GetDeviceByID(ctx, resource.DeviceID)
	if err == nil {
		_ = d.Set("serial_number", device.SerialNumber)
	}
	return diags
}

func resourceSTLAppCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	var diags diag.Diagnostics
	var client *stl.Client
	var err error

	if endpoint, ok := d.GetOk("endpoint"); ok {
		client, err = config.STLClient(endpoint.(string))
	} else {
		client, err = config.STLClient()
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
		return diag.FromErr(fmt.Errorf("stl_app: create STL app: %w", err))
	}
	d.SetId(fmt.Sprintf("%d", resource.ID))
	syncSTLIfNeeded(ctx, client, d, m)
	return diags
}
