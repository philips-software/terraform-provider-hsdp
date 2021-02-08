package hsdp

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/stl"
)

func resourceSTLConfig() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceSTLDeviceConfigCreate,
		ReadContext:   resourceSTLDeviceConfigRead,
		UpdateContext: resourceSTLDeviceConfigUpdate,
		DeleteContext: resourceSTLDeviceConfigDelete,

		Schema: map[string]*schema.Schema{
			"serial_number": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"firewall_exceptions": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tcp": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"udp": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"logging": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"raw_config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hsdp_product_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hsdp_shared_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hsdp_secret_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hsdp_ingestor_host": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hsdp_custom_field": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"cert": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 100,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
		},
	}
}

func resourceDataToInput(fwExceptions *stl.AppFirewallException, appLogging *stl.AppLogging, d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	if d == nil {
		return fmt.Errorf("dataToResourceData: schema.ResourceData is nil")
	}
	// Logging
	if v, ok := d.GetOk("logging"); ok {
		vL := v.(*schema.Set).List()
		for i, vi := range vL {
			config.Debug("Reading Logging Set %d\n", i)
			mVi := vi.(map[string]interface{})
			fwExceptions.TCP = expandIntList(mVi["tcp"].(*schema.Set).List())
			fwExceptions.UDP = expandIntList(mVi["udp"].(*schema.Set).List())
		}
	}
	// Firewall exceptions
	if v, ok := d.GetOk("firewall_exceptions"); ok {
		vL := v.(*schema.Set).List()
		for i, vi := range vL {
			config.Debug("Reading Firewall exception Set %d\n", i)
			mVi := vi.(map[string]interface{})
			appLogging.HSDPIngestorHost = mVi["hsdp_ingestor_host"].(string)
			appLogging.HSDPProductKey = mVi["hsdp_product_key"].(string)
			appLogging.HSDPSharedKey = mVi["hsdp_shared_key"].(string)
			appLogging.HSDPSecretKey = mVi["hsdp_secret_key"].(string)
			appLogging.HSDPCustomField = mVi["hsdp_custom_field"].(bool)
			appLogging.RawConfig = mVi["raw_config"].(string)
		}
	}
	// TODO: certs

	return nil
}

func dataToResourceData(fwExceptions *stl.AppFirewallException, appLogging *stl.AppLogging, d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	if d == nil {
		return fmt.Errorf("dataToResourceData: schema.ResourceData is nil")
	}
	// Logging
	s := &schema.Set{F: resourceMetricsThresholdHash}
	appLoggingDef := make(map[string]interface{})
	appLoggingDef["raw_config"] = appLogging.RawConfig
	appLoggingDef["hsdp_product_key"] = appLogging.HSDPProductKey
	appLoggingDef["hsdp_shared_key"] = appLogging.HSDPSharedKey
	appLoggingDef["hsdp_secret_key"] = appLogging.HSDPSecretKey
	appLoggingDef["hsdp_ingestor_host"] = appLogging.HSDPIngestorHost
	appLoggingDef["hsdp_custom_field"] = appLogging.HSDPCustomField
	s.Add(appLoggingDef)
	_, _ = config.Debug("Adding logging data")
	err := d.Set("logging", s)
	if err != nil {
		return fmt.Errorf("dataToResourceData: logging: %w", err)
	}
	// Firewall exceptions
	s = &schema.Set{F: resourceMetricsThresholdHash}
	fwExceptionsDef := make(map[string]interface{})
	fwExceptionsDef["tcp"] = fwExceptions.TCP
	fwExceptionsDef["udp"] = fwExceptions.UDP
	s.Add(fwExceptionsDef)
	_, _ = config.Debug("Adding firewall exceptions data")
	err = d.Set("firewall_exceptions", s)
	if err != nil {
		return fmt.Errorf("dataToResourceData: firewall_exceptions: %w", err)
	}
	// TODO: certs

	return nil
}

func resourceSTLDeviceConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceSTLDeviceConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceSTLDeviceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	serialNumber := d.Id()
	firewallExceptions, err := client.Config.GetFirewallExceptionsBySerial(ctx, serialNumber)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read firewall exceptions: %w", err))
	}
	appLogging, err := client.Config.GetAppLoggingBySerial(ctx, serialNumber)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read appLogging: %w", err))
	}
	err = dataToResourceData(firewallExceptions, appLogging, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceSTLDeviceConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
