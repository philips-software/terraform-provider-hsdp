package mdm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dip-software/go-dip-api/connect/mdm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceConnectMDMFirmwareDistributionRequest() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext:      resourceConnectMDMFirmwareDistributionRequestCreate,
		ReadContext:        resourceConnectMDMFirmwareDistributionRequestRead,
		UpdateContext:      resourceConnectMDMFirmwareDistributionRequestUpdate,
		DeleteContext:      resourceConnectMDMFirmwareDistributionRequestDelete,
		DeprecationMessage: "This will be replace by the Firmware v2 API. Only use it for test/demo purposes!",

		Schema: map[string]*schema.Schema{
			"firmware_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"orchestration_mode": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"user_consent_required": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"distribution_target_device_groups_ids": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				MaxItems: 10,
				Elem:     tools.StringSchema(),
			},
			"firmware_component_version_ids": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				MaxItems: 5,
				Elem:     tools.StringSchema(),
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

func schemaToFirmwareDistributionRequest(d *schema.ResourceData) mdm.FirmwareDistributionRequest {
	firmwareVersion := d.Get("firmware_version").(string)
	description := d.Get("description").(string)
	orchestrationMode := d.Get("orchestration_mode").(string)
	userConsentRequired := d.Get("user_consent_required").(bool)
	status := d.Get("status").(string)

	targetDeviceGroupIds := tools.ExpandStringList(d.Get("distribution_target_device_groups_ids").(*schema.Set).List())
	componentVersionIds := tools.ExpandStringList(d.Get("firmware_component_version_ids").(*schema.Set).List())

	resource := mdm.FirmwareDistributionRequest{
		FirmwareVersion:     firmwareVersion,
		Description:         description,
		OrchestrationMode:   orchestrationMode,
		UserConsentRequired: userConsentRequired,
		Status:              status,
	}
	var dgs []mdm.Reference
	for _, g := range targetDeviceGroupIds {
		dgs = append(dgs, mdm.Reference{
			Reference: g,
		})
	}
	resource.DistributionTargets = dgs

	var compVers []mdm.Reference
	for _, v := range componentVersionIds {
		compVers = append(compVers, mdm.Reference{
			Reference: v,
		})
	}
	resource.FirmwareComponentVersions = compVers

	return resource
}

func firmwareDistributionRequestToSchema(resource mdm.FirmwareDistributionRequest, d *schema.ResourceData) {
	_ = d.Set("description", resource.Description)
	_ = d.Set("firmware_version", resource.FirmwareVersion)
	_ = d.Set("orchestration_mode", resource.OrchestrationMode)
	_ = d.Set("user_consent_required", resource.UserConsentRequired)
	_ = d.Set("guid", resource.ID)
	_ = d.Set("status", resource.Status)

	var dts []string
	for _, g := range resource.DistributionTargets {
		dts = append(dts, g.Reference)
	}
	_ = d.Set("distribution_target_device_groups_ids", tools.SchemaSetStrings(dts))

	var cvs []string
	for _, v := range resource.FirmwareComponentVersions {
		cvs = append(cvs, v.Reference)
	}
	_ = d.Set("firmware_component_version_ids", tools.SchemaSetStrings(cvs))
}

func resourceConnectMDMFirmwareDistributionRequestCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToFirmwareDistributionRequest(d)

	var created *mdm.FirmwareDistributionRequest
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.FirmwareDistributionRequests.Create(resource)
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
		return diag.FromErr(fmt.Errorf("failed to create resource: %d", resp.StatusCode()))
	}
	_ = d.Set("guid", created.ID)
	d.SetId(fmt.Sprintf("FirmwareDistributionRequest/%s", created.ID))

	return resourceConnectMDMFirmwareDistributionRequestRead(ctx, d, m)
}

func resourceConnectMDMFirmwareDistributionRequestRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "FirmwareDistributionRequest/%s", &id)
	var resource *mdm.FirmwareDistributionRequest
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.FirmwareDistributionRequests.GetByID(id)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp != nil && (resp.StatusCode() == http.StatusNotFound || resp.StatusCode() == http.StatusGone) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	firmwareDistributionRequestToSchema(*resource, d)
	return diags
}

func resourceConnectMDMFirmwareDistributionRequestUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)

	resource := schemaToFirmwareDistributionRequest(d)
	resource.ID = id

	if !d.HasChange("status") {
		return diag.FromErr(fmt.Errorf("only the 'status' can be updated"))
	}

	_, _, err = client.FirmwareDistributionRequests.Update(resource)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		return diags
	}
	return resourceConnectMDMFirmwareDistributionRequestRead(ctx, d, m)
}

func resourceConnectMDMFirmwareDistributionRequestDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	resource, _, err := client.FirmwareDistributionRequests.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.FirmwareDistributionRequests.Delete(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
