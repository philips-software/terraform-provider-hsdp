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
		CreateContext: resourceSTLConfigCreate,
		ReadContext:   resourceSTLConfigRead,
		UpdateContext: resourceSTLConfigUpdate,
		DeleteContext: resourceSTLConfigDelete,

		Schema: map[string]*schema.Schema{
			"serial_number": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sync": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
							MaxItems: 65535,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"udp": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 65535,
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
						"hsdp_logging": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
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
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func resourceSTLConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	loggingRef := stl.UpdateAppLoggingInput{
		SerialNumber: serialNumber,
	}
	fwExceptionRef := stl.UpdateAppFirewallExceptionInput{
		SerialNumber: serialNumber,
	}
	fwExceptionRef.TCP = []int{}
	fwExceptionRef.UDP = []int{}
	// Clear
	if _, ok := d.GetOk("logging"); ok {
		_, err = client.Config.UpdateAppLogging(ctx, loggingRef)
		if err != nil {
			return diag.FromErr(fmt.Errorf("hsdp_stl_config: UpdateAppLogging: %w", err))
		}
	}
	if _, ok := d.GetOk("firewall_exceptions"); ok {
		_, err = client.Config.UpdateAppFirewallExceptions(ctx, fwExceptionRef)
		if err != nil {
			return diag.FromErr(fmt.Errorf("hsdp_stl_config: UpdateAppFirewallExceptions: %w", err))
		}
	}
	syncSTLIfNeeded(ctx, client, d, m)
	d.SetId("")
	return diags
}

func resourceSTLConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceSTLConfigCreate(ctx, d, m)
}

func resourceDataToInput(fwExceptions *stl.UpdateAppFirewallExceptionInput, appLogging *stl.UpdateAppLoggingInput, d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	if d == nil {
		return fmt.Errorf("dataToResourceData: schema.ResourceData is nil")
	}
	serialNumber := d.Get("serial_number").(string)
	// check if serialNumber checks out, if not we may need to fetch by ID

	// Firewall exceptions
	if v, ok := d.GetOk("firewall_exceptions"); ok {
		vL := v.(*schema.Set).List()
		for i, vi := range vL {
			_, _ = config.Debug("Reading Logging Set %d\n", i)
			mVi := vi.(map[string]interface{})
			fwExceptions.TCP = expandIntList(mVi["tcp"].(*schema.Set).List())
			fwExceptions.UDP = expandIntList(mVi["udp"].(*schema.Set).List())
		}
	}
	fwExceptions.SerialNumber = serialNumber

	// Logging
	if v, ok := d.GetOk("logging"); ok {
		vL := v.(*schema.Set).List()
		for i, vi := range vL {
			_, _ = config.Debug("Reading Firewall exception Set %d\n", i)
			mVi := vi.(map[string]interface{})
			if a, ok := mVi["hsdp_logging"].(bool); ok {
				appLogging.HSDPLogging = a
			}
			if a, ok := mVi["hsdp_ingestor_host"].(string); ok {
				appLogging.HSDPIngestorHost = a
			}
			if a, ok := mVi["hsdp_product_key"].(string); ok {
				appLogging.HSDPProductKey = a
			}
			if a, ok := mVi["hsdp_shared_key"].(string); ok {
				appLogging.HSDPSharedKey = a
			}
			if a, ok := mVi["hsdp_secret_key"].(string); ok {
				appLogging.HSDPSecretKey = a
			}
			if a, ok := mVi["hsdp_custom_field"].(bool); ok {
				appLogging.HSDPCustomField = &a
			}
			if a, ok := mVi["raw_config"].(string); ok {
				appLogging.RawConfig = a
			}
		}
		if ok, err := appLogging.Validate(); !ok {
			return err
		}
	}
	appLogging.SerialNumber = serialNumber

	return nil
}

func dataToResourceData(fwExceptions *stl.AppFirewallException, appLogging *stl.AppLogging, d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	if d == nil {
		return fmt.Errorf("dataToResourceData: schema.ResourceData is nil")
	}
	// Logging
	if _, ok := d.GetOk("logging"); ok {
		s := &schema.Set{F: resourceMetricsThresholdHash}
		appLoggingDef := make(map[string]interface{})
		appLoggingDef["raw_config"] = appLogging.RawConfig
		appLoggingDef["hsdp_product_key"] = appLogging.HSDPProductKey
		appLoggingDef["hsdp_shared_key"] = appLogging.HSDPSharedKey
		appLoggingDef["hsdp_secret_key"] = appLogging.HSDPSecretKey
		appLoggingDef["hsdp_ingestor_host"] = appLogging.HSDPIngestorHost
		appLoggingDef["hsdp_custom_field"] = appLogging.HSDPCustomField
		appLoggingDef["hsdp_logging"] = appLogging.HSDPLogging
		s.Add(appLoggingDef)
		_, _ = config.Debug("Adding logging data")
		err := d.Set("logging", s)
		if err != nil {
			return fmt.Errorf("dataToResourceData: logging: %w", err)
		}
	}
	// Firewall exceptions
	if _, ok := d.GetOk("firewall_exceptions"); ok {
		s := &schema.Set{F: resourceMetricsThresholdHash}
		fwExceptionsDef := make(map[string]interface{})
		fwExceptionsDef["tcp"] = fwExceptions.TCP
		fwExceptionsDef["udp"] = fwExceptions.UDP
		s.Add(fwExceptionsDef)
		_, _ = config.Debug("Adding firewall exceptions data")
		err := d.Set("firewall_exceptions", s)
		if err != nil {
			return fmt.Errorf("dataToResourceData: firewall_exceptions: %w", err)
		}
	}
	return nil
}

func resourceSTLConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceSTLConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	var loggingRef stl.UpdateAppLoggingInput
	var fwExceptionRef stl.UpdateAppFirewallExceptionInput
	err = resourceDataToInput(&fwExceptionRef, &loggingRef, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	if _, ok := d.GetOk("logging"); ok {
		_, err = client.Config.UpdateAppLogging(ctx, loggingRef)
		if err != nil {
			return diag.FromErr(fmt.Errorf("hsdp_stl_config: UpdateAppLogging: %w", err))
		}
	}
	if _, ok := d.GetOk("firewall_exceptions"); ok {
		_, err = client.Config.UpdateAppFirewallExceptions(ctx, fwExceptionRef)
		if err != nil {
			return diag.FromErr(fmt.Errorf("hsdp_stl_config: UpdateAppFirewallExceptions: %w", err))
		}
	}
	if d.IsNewResource() {
		d.SetId(loggingRef.SerialNumber)
	}
	syncSTLIfNeeded(ctx, client, d, m)
	return diags
}
