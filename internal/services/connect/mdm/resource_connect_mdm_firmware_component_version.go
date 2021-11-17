package mdm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceConnectMDMFirmwareComponentVersion() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceConnectMDMFirmwareComponentVersionCreate,
		ReadContext:   resourceConnectMDMFirmwareComponentVersionRead,
		UpdateContext: resourceConnectMDMFirmwareComponentVersionUpdate,
		DeleteContext: resourceConnectMDMFirmwareComponentVersionDelete,

		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"effective_date": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"firmware_component_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"blob_url": {
				Type:     schema.TypeString,
				Default:  true,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Default:  true,
				Optional: true,
			},
			"fingerprint": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem:     fingerprintSchema(),
			},
			"component_required": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Required: true,
			},
			"custom_resource": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deprecated_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func fingerprintSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"algorithm": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hash": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func schemaToFirmwareComponentVersion(d *schema.ResourceData) mdm.FirmwareComponentVersion {
	version := d.Get("version").(string)
	description := d.Get("description").(string)
	firmwareComponentId := d.Get("firmware_component_id").(string)
	effectiveDate := d.Get("effective_date").(string)
	deprecatedDate := d.Get("deprecated_date").(string)
	componentRequired := d.Get("component_required").(bool)
	customResource := d.Get("custom_resource").(string)
	blobUrl := d.Get("blob_url").(string)
	size := d.Get("size").(int)

	resource := mdm.FirmwareComponentVersion{
		Version:           version,
		Description:       description,
		EffectiveDate:     effectiveDate,
		DeprecatedDate:    deprecatedDate,
		ComponentRequired: componentRequired,
		CustomResource:    []byte(customResource),
		BlobURL:           blobUrl,
		Size:              size,
	}

	if v, ok := d.GetOk("fingerprint"); ok {
		vL := v.(*schema.Set).List()
		for _, entry := range vL {
			mV := entry.(map[string]interface{})
			resource.FingerPrint = &mdm.Fingerprint{
				Algorithm: mV["algorithm"].(string),
				Hash:      mV["hash"].(string),
			}
		}
	}
	resource.FirmwareComponentId.Reference = firmwareComponentId
	// TODO: add encryption info

	return resource
}

func FirmwareComponentVersionToSchema(resource mdm.FirmwareComponentVersion, d *schema.ResourceData) {
	_ = d.Set("version", resource.Version)
	_ = d.Set("description", resource.Description)
	_ = d.Set("firmware_component_id", resource.FirmwareComponentId.Reference)
	_ = d.Set("custom_resource", string(resource.CustomResource))
	_ = d.Set("effective_date", resource.EffectiveDate)
	_ = d.Set("deprecated_date", resource.DeprecatedDate)
	_ = d.Set("blob_url", resource.BlobURL)
	_ = d.Set("size", resource.Size)
	_ = d.Set("component_required", resource.ComponentRequired)
	_ = d.Set("guid", resource.ID)

	// Add Fingerprint
	a := &schema.Set{F: schema.HashResource(fingerprintSchema())}
	entry := make(map[string]interface{})
	entry["algorithm"] = resource.FingerPrint.Algorithm
	entry["hash"] = resource.FingerPrint.Hash
	a.Add(entry)
	_ = d.Set("fingerprint", a)
}

func resourceConnectMDMFirmwareComponentVersionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToFirmwareComponentVersion(d)

	var created *mdm.FirmwareComponentVersion
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.FirmwareComponentVersions.Create(resource)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if created == nil {
		return diag.FromErr(fmt.Errorf("failed to create resource: %d", resp.StatusCode))
	}
	_ = d.Set("guid", created.ID)
	d.SetId(fmt.Sprintf("FirmwareComponentVersion/%s", created.ID))
	return resourceConnectMDMFirmwareComponentVersionRead(ctx, d, m)
}

func resourceConnectMDMFirmwareComponentVersionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "FirmwareComponentVersion/%s", &id)
	service, resp, err := client.FirmwareComponentVersions.GetByID(id)
	if err != nil {
		if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	FirmwareComponentVersionToSchema(*service, d)
	return diags
}

func resourceConnectMDMFirmwareComponentVersionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)

	service := schemaToFirmwareComponentVersion(d)
	service.ID = id

	_, _, err = client.FirmwareComponentVersions.Update(service)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		return diags
	}
	return resourceConnectMDMFirmwareComponentVersionRead(ctx, d, m)
}

func resourceConnectMDMFirmwareComponentVersionDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	resource, _, err := client.FirmwareComponentVersions.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.FirmwareComponentVersions.Delete(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
