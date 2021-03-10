package hsdp

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/stl"
	"log"
	"sort"
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
							Optional: true,
							MaxItems: 1024,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"udp": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1024,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"clear_on_destroy": {
							Description: "Clear all exceptions on resource destroy",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"ensure_tcp": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1024,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"ensure_udp": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1024,
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

func containsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
	if clearFirewallExceptionsOnDestroy(d) {
		_, err = client.Config.UpdateAppFirewallExceptions(ctx, fwExceptionRef)
	}
	syncSTLIfNeeded(ctx, client, d, m)
	d.SetId("")
	return diags
}

func resourceSTLConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceSTLConfigCreate(ctx, d, m)
}

func resourceDataToInput(ctx context.Context, client *stl.Client, fwExceptions *stl.UpdateAppFirewallExceptionInput, appLogging *stl.UpdateAppLoggingInput, d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	if d == nil {
		return fmt.Errorf("dataToResourceData: schema.ResourceData is nil")
	}
	serialNumber := d.Get("serial_number").(string)
	// check if serialNumber checks out, if not we may need to fetch by ID

	// Firewall exceptions
	if err := validateFirewallExceptions(d); err != nil {
		return err
	}
	tcp := []int{}
	udp := []int{}
	ensureTCP := []int{}
	ensureUDP := []int{}

	if v, ok := d.GetOk("firewall_exceptions"); ok {
		vL := v.(*schema.Set).List()
		for i, vi := range vL {
			_, _ = config.Debug("Reading Logging Set %d\n", i)
			mVi := vi.(map[string]interface{})
			if _, ok := mVi["tcp"].(*schema.Set); ok {
				tcp = expandIntList(mVi["tcp"].(*schema.Set).List())
			}
			if _, ok := mVi["udp"].(*schema.Set); ok {
				udp = expandIntList(mVi["udp"].(*schema.Set).List())
			}
			ensureTCP = expandIntList(mVi["ensure_tcp"].(*schema.Set).List())
			ensureUDP = expandIntList(mVi["ensure_udp"].(*schema.Set).List())
		}
	}
	if len(tcp) > 0 {
		fwExceptions.TCP = tcp
	}
	if len(udp) > 0 {
		fwExceptions.UDP = udp
	}
	if len(ensureTCP) > 0 || len(ensureUDP) > 0 { // Fetch current settings
		currentSettings, err := client.Config.GetFirewallExceptionsBySerial(ctx, serialNumber)
		if err != nil {
			return err
		}
		if len(ensureTCP) > 0 {
			fwExceptions.TCP = mergePorts(currentSettings.TCP, ensureTCP)
		}
		if len(ensureUDP) > 0 {
			fwExceptions.UDP = mergePorts(currentSettings.UDP, ensureUDP)
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

func mergePorts(i []int, ensure []int) []int {
	// Sort
	ports := append(i, ensure...)
	sort.Ints(ports)
	// Unique
	seen := make(map[int]struct{}, len(ports))
	j := 0
	for _, v := range ports {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		ports[j] = v
		j++
	}
	return ports[:j]
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
		ensureTCP := getPortList(d, "ensure_tcp")
		actualTCP := []int{}
		for _, p := range ensureTCP {
			if containsInt(fwExceptions.TCP, p) {
				actualTCP = append(actualTCP, p)
			}
		}
		fwExceptionsDef["ensure_tcp"] = actualTCP
		ensureUDP := getPortList(d, "ensure_udp")
		actualUDP := []int{}
		for _, p := range ensureUDP {
			if containsInt(fwExceptions.UDP, p) {
				actualUDP = append(actualUDP, p)
			}
		}
		fwExceptionsDef["ensure_udp"] = actualUDP
		if !hasFirewallExceptionField(d, "ensure_tcp") {
			fwExceptionsDef["tcp"] = fwExceptions.TCP
		}
		if !hasFirewallExceptionField(d, "ensure_udp") {
			fwExceptionsDef["udp"] = fwExceptions.UDP
		}
		s.Add(fwExceptionsDef)
		_, _ = config.Debug("Adding firewall exceptions data\n")
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
	err = resourceDataToInput(ctx, client, &fwExceptionRef, &loggingRef, d, m)
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
	return resourceSTLConfigRead(ctx, d, m)
}

func clearFirewallExceptionsOnDestroy(d *schema.ResourceData) bool {
	if v, ok := d.GetOk("firewall_exceptions"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			if choice, ok := mVi["clear_on_destroy"].(bool); ok {
				return choice
			}
		}
	}
	return false
}

func hasFirewallExceptionField(d *schema.ResourceData, fieldName string) bool {
	if v, ok := d.GetOk("firewall_exceptions"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			if mVi[fieldName] != nil {
				return true
			}
		}
	}
	return false
}

func getPortList(d *schema.ResourceData, field string) []int {
	if v, ok := d.GetOk("firewall_exceptions"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			portList := expandIntList(mVi[field].(*schema.Set).List())
			sort.Ints(portList)
			return portList
		}
	}
	return []int{}
}

func validateFirewallExceptions(d *schema.ResourceData) error {
	log.Printf("Validating firewall Exceptions\n")
	if v, ok := d.GetOk("firewall_exceptions"); ok {
		foundTCP := []int{}
		foundUDP := []int{}
		foundEnsureTCP := []int{}
		foundEnsureUDP := []int{}
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			foundTCP = expandIntList(mVi["tcp"].(*schema.Set).List())
			foundUDP = expandIntList(mVi["udp"].(*schema.Set).List())
			foundEnsureTCP = expandIntList(mVi["ensure_tcp"].(*schema.Set).List())
			foundEnsureUDP = expandIntList(mVi["ensure_udp"].(*schema.Set).List())
		}
		if len(foundEnsureTCP) > 0 && len(foundTCP) > 0 {
			return fmt.Errorf("conflicting 'ensure_tcp' and 'tcp")
		}
		if len(foundEnsureUDP) > 0 && len(foundUDP) > 0 {
			return fmt.Errorf("conflicting 'ensure_udp' and 'udp")
		}
		if len(foundTCP) > 0 && len(foundEnsureTCP) > 0 {
			return fmt.Errorf("conflicting 'tcp' and 'ensure_tcp")
		}
		if len(foundUDP) > 0 && len(foundEnsureUDP) > 0 {
			return fmt.Errorf("conflicting 'udp' and 'ensure_udp")
		}
	}
	return nil
}
