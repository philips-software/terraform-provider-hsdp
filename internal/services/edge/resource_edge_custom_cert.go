package edge

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/stl"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceEdgeCustomCert() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceEdgeCustomCertCreate,
		ReadContext:   resourceEdgeCustomCertRead,
		UpdateContext: resourceEdgeCustomCertUpdate,
		DeleteContext: resourceEdgeCustomCertDelete,

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
			"private_key_pem": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cert_pem": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sync": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceEdgeCustomCertDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	_, err = client.Certs.DeleteCustomCert(ctx, stl.DeleteAppCustomCertInput{ID: resourceID})
	if err != nil {
		return diag.FromErr(fmt.Errorf("stl_custom_cert delete: %w", err))
	}
	syncSTLIfNeeded(ctx, client, d, m)
	d.SetId("")
	return diags
}

func resourceEdgeCustomCertUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	_, err = client.Certs.UpdateCustomCert(ctx, stl.UpdateAppCustomCertInput{
		ID:   resourceID,
		Name: d.Get("name").(string),
		Key:  d.Get("private_key_pem").(string),
		Cert: d.Get("cert_pem").(string),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("stl_custom_cert update: %w", err))
	}
	syncSTLIfNeeded(ctx, client, d, m)
	return diags
}

func resourceEdgeCustomCertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	readCert, err := client.Certs.GetCustomCertByID(ctx, resourceID)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("name", readCert.Name)
	_ = d.Set("cert_pem", readCert.Cert)
	_ = d.Set("private_key_pem", readCert.Key)
	return diags
}

func resourceEdgeCustomCertCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	serialNumber := d.Get("serial_number").(string)
	newCert := stl.CreateAppCustomCertInput{
		SerialNumber: serialNumber,
	}
	newCert.Name = d.Get("name").(string)
	newCert.Key = d.Get("private_key_pem").(string)
	newCert.Cert = d.Get("cert_pem").(string)
	created, err := client.Certs.CreateCustomCert(ctx, newCert)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%d", created.ID))
	syncSTLIfNeeded(ctx, client, d, m)
	return diags
}

func syncSTLIfNeeded(ctx context.Context, client *stl.Client, d *schema.ResourceData, m interface{}) {
	c := m.(*config.Config)
	sync := d.Get("sync").(bool)
	if !sync {
		return
	}
	serialNumber := d.Get("serial_number").(string)
	_, _ = c.Debug("Syncing %s\n", serialNumber)
	_ = client.Devices.SyncDeviceConfig(ctx, serialNumber)
}
