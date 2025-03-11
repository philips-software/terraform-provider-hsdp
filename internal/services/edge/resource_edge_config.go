package edge

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/dip-software/go-dip-api/stl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

const (
	clearOnDestroyDefault = true
)

func ResourceEdgeConfig() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceEdgeConfigCreate,
		ReadContext:   resourceEdgeConfigRead,
		UpdateContext: resourceEdgeConfigUpdate,
		DeleteContext: resourceEdgeConfigDelete,

		Schema: map[string]*schema.Schema{
			"serial_number": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"principal": config.PrincipalSchema(),
			"sync": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"firewall_exceptions": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem:     firewallExceptionsSchema(),
			},
			"logging": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem:     loggingElem(),
			},
		},
	}
}

func loggingElem() *schema.Resource {
	return &schema.Resource{
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
	}
}
func firewallExceptionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"tcp": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1024,
				Elem:     tools.IntSchema(),
			},
			"udp": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1024,
				Elem:     tools.IntSchema(),
			},
			"clear_on_destroy": {
				Description: "Clear ports on resource destroy",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     clearOnDestroyDefault,
			},
			"ensure_tcp": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1024,
				Elem:     tools.IntSchema(),
			},
			"ensure_udp": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1024,
				Elem:     tools.IntSchema(),
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

func resourceEdgeConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	loggingRef := stl.UpdateAppLoggingInput{
		SerialNumber: serialNumber,
	}
	fwExceptionRef := stl.UpdateAppFirewallExceptionInput{
		SerialNumber: serialNumber,
	}
	// Clear
	if _, ok := d.GetOk("logging"); ok {
		_, err = client.Config.UpdateAppLogging(ctx, loggingRef)
		if err != nil {
			return diag.FromErr(fmt.Errorf("hsdp_edge_config: UpdateAppLogging: %w", err))
		}
	}
	if _, ok := d.GetOk("firewall_exceptions"); ok && clearFirewallExceptionsOnDestroy(d) {
		currentSettings, err := client.Config.GetFirewallExceptionsBySerial(ctx, serialNumber)
		if err != nil {
			return diag.FromErr(fmt.Errorf("delete Edge discovery: %w", err))
		}
		fwExceptionRef.TCP = currentSettings.TCP
		fwExceptionRef.UDP = currentSettings.UDP
		if hasFirewallExceptionField(d, "tcp") {
			fwExceptionRef.TCP = []int{}
		}
		if hasFirewallExceptionField(d, "udp") {
			fwExceptionRef.UDP = []int{}
		}
		if hasFirewallExceptionField(d, "ensure_tcp") {
			pruneList := getPortList(d, "ensure_tcp")
			fwExceptionRef.TCP = tools.PrunePorts(currentSettings.TCP, pruneList)
		}
		if hasFirewallExceptionField(d, "ensure_udp") {
			pruneList := getPortList(d, "ensure_udp")
			fwExceptionRef.UDP = tools.PrunePorts(currentSettings.UDP, pruneList)
		}
		_, err = client.Config.UpdateAppFirewallExceptions(ctx, fwExceptionRef)
		if err != nil {
			return diag.FromErr(fmt.Errorf("hsdp_edge_config: UpdateAppFirewallExceptions: %w", err))
		}
	}
	syncSTLIfNeeded(ctx, client, d, m)
	d.SetId("")
	return diags
}

func resourceEdgeConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceEdgeConfigCreate(ctx, d, m)
}

func resourceDataToInput(ctx context.Context, client *stl.Client, fwExceptions *stl.UpdateAppFirewallExceptionInput, appLogging *stl.UpdateAppLoggingInput, d *schema.ResourceData, m interface{}) error {
	c := m.(*config.Config)

	if d == nil {
		return fmt.Errorf("dataToResourceData: schema.ResourceData is nil")
	}
	serialNumber := d.Get("serial_number").(string)
	// check if serialNumber checks out, if not we may need to fetch by ID

	// Firewall exceptions
	if err := validateFirewallExceptions(d); err != nil {
		return err
	}
	var tcp []int
	var udp []int
	var ensureTCP []int
	var ensureUDP []int

	if v, ok := d.GetOk("firewall_exceptions"); ok {
		vL := v.(*schema.Set).List()
		for i, vi := range vL {
			_, _ = c.Debug("Reading Logging Set %d\n", i)
			mVi := vi.(map[string]interface{})
			if _, ok := mVi["tcp"].(*schema.Set); ok {
				tcp = tools.ExpandIntList(mVi["tcp"].(*schema.Set).List())
			}
			if _, ok := mVi["udp"].(*schema.Set); ok {
				udp = tools.ExpandIntList(mVi["udp"].(*schema.Set).List())
			}
			ensureTCP = tools.ExpandIntList(mVi["ensure_tcp"].(*schema.Set).List())
			ensureUDP = tools.ExpandIntList(mVi["ensure_udp"].(*schema.Set).List())
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
			_, _ = c.Debug("Reading Firewall exception Set %d\n", i)
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
	c := m.(*config.Config)

	if d == nil {
		return fmt.Errorf("dataToResourceData: schema.ResourceData is nil")
	}
	// Logging
	if _, ok := d.GetOk("logging"); ok {
		s := &schema.Set{F: schema.HashResource(loggingElem())}
		appLoggingDef := make(map[string]interface{})
		appLoggingDef["raw_config"] = appLogging.RawConfig
		appLoggingDef["hsdp_product_key"] = appLogging.HSDPProductKey
		appLoggingDef["hsdp_shared_key"] = appLogging.HSDPSharedKey
		appLoggingDef["hsdp_secret_key"] = appLogging.HSDPSecretKey
		appLoggingDef["hsdp_ingestor_host"] = appLogging.HSDPIngestorHost
		appLoggingDef["hsdp_custom_field"] = appLogging.HSDPCustomField
		appLoggingDef["hsdp_logging"] = appLogging.HSDPLogging
		s.Add(appLoggingDef)
		_, _ = c.Debug("Adding logging data")
		err := d.Set("logging", s)
		if err != nil {
			return fmt.Errorf("dataToResourceData: logging: %w", err)
		}
	}
	// Firewall exceptions
	if _, ok := d.GetOk("firewall_exceptions"); ok {
		s := &schema.Set{F: schema.HashResource(firewallExceptionsSchema())}
		fwExceptionsDef := make(map[string]interface{})

		// Determine actual TCP/UDP settings
		ensureTCP := getPortList(d, "ensure_tcp")
		var actualTCP []int
		for _, p := range ensureTCP {
			if containsInt(fwExceptions.TCP, p) {
				actualTCP = append(actualTCP, p)
			}
		}
		fwExceptionsDef["ensure_tcp"] = tools.SchemaSetInts(actualTCP)
		ensureUDP := getPortList(d, "ensure_udp")
		var actualUDP []int
		for _, p := range ensureUDP {
			if containsInt(fwExceptions.UDP, p) {
				actualUDP = append(actualUDP, p)
			}
		}
		fwExceptionsDef["ensure_udp"] = tools.SchemaSetInts(actualUDP)

		// This is explicit TCP/UDP port list mode
		if !hasFirewallExceptionField(d, "ensure_tcp") {
			fwExceptionsDef["tcp"] = tools.SchemaSetInts(fwExceptions.TCP)
		}
		if !hasFirewallExceptionField(d, "ensure_udp") {
			fwExceptionsDef["udp"] = tools.SchemaSetInts(fwExceptions.UDP)
		}
		fwExceptionsDef["clear_on_destroy"] = clearFirewallExceptionsOnDestroy(d)
		s.Add(fwExceptionsDef)
		_, _ = c.Debug("Adding firewall exceptions data\n")
		err := d.Set("firewall_exceptions", s)
		if err != nil {
			return fmt.Errorf("dataToResourceData: firewall_exceptions: %w", err)
		}
	}
	return nil
}

func resourceEdgeConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics
	var client *stl.Client
	var err error

	principal := config.SchemaToPrincipal(d, m)
	client, err = c.STLClient(principal)

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

func resourceEdgeConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var client *stl.Client
	var err error

	principal := config.SchemaToPrincipal(d, m)
	client, err = c.STLClient(principal)

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
			return diag.FromErr(fmt.Errorf("hsdp_edge_config: UpdateAppLogging: %w", err))
		}
	}
	if _, ok := d.GetOk("firewall_exceptions"); ok {
		_, err = client.Config.UpdateAppFirewallExceptions(ctx, fwExceptionRef)
		if err != nil {
			return diag.FromErr(fmt.Errorf("hsdp_edge_config: UpdateAppFirewallExceptions: %w", err))
		}
	}
	if d.IsNewResource() {
		d.SetId(loggingRef.SerialNumber)
	}
	syncSTLIfNeeded(ctx, client, d, m)
	return resourceEdgeConfigRead(ctx, d, m)
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
	return clearOnDestroyDefault
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

func getPortList(d *schema.ResourceData, field string) []int {
	if v, ok := d.GetOk("firewall_exceptions"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			portList := tools.ExpandIntList(mVi[field].(*schema.Set).List())
			sort.Ints(portList)
			return portList
		}
	}
	return []int{}
}

func validateFirewallExceptions(d *schema.ResourceData) error {
	log.Printf("Validating firewall Exceptions\n")
	if v, ok := d.GetOk("firewall_exceptions"); ok {
		var foundTCP []int
		var foundUDP []int
		var foundEnsureTCP []int
		var foundEnsureUDP []int
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			foundTCP = tools.ExpandIntList(mVi["tcp"].(*schema.Set).List())
			foundUDP = tools.ExpandIntList(mVi["udp"].(*schema.Set).List())
			foundEnsureTCP = tools.ExpandIntList(mVi["ensure_tcp"].(*schema.Set).List())
			foundEnsureUDP = tools.ExpandIntList(mVi["ensure_udp"].(*schema.Set).List())
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
