package device

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dip-software/go-dip-api/iam"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

var descriptions = map[string]string{
	"device": "These resources represent device accounts in IAM. Typically, devices contain information that identifies a device’s uniqueness, intended use, credentials, and other details to track the device back to its proposition",
}

func ResourceIAMDevice() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceIAMDeviceCreate,
		ReadContext:   resourceIAMDeviceRead,
		UpdateContext: resourceIAMDeviceUpdate,
		DeleteContext: resourceIAMDeviceDelete,
		Description:   descriptions["device"],

		Schema: map[string]*schema.Schema{
			"login_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The login id of the device.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The password of the device.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The application ID (GUID) this device should be attached to.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The organization ID (GUID) this device should be attached to.",
			},
			"registration_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the device was registered.",
			},
			"debug_until": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"text": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"for_test": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "When set to true this device is marked as a test device.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Controls if this device is active or not.",
			},
			"external_identifier": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Elem:        externalIdentifierSchema(),
				Description: "Block describing external ID of this device.",
			},
			"global_reference_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressWhenGenerated,
			},
		},
	}
}

func externalIdentifierSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"system": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     deviceTypeSchema(),
			},
		},
	}
}

func deviceTypeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"text": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func schemaToDevice(d *schema.ResourceData) iam.Device {
	loginId := d.Get("login_id").(string)
	password := d.Get("password").(string)
	//text := d.Get("text").(string)
	organizationId := d.Get("organization_id").(string)
	applicationId := d.Get("application_id").(string)
	deviceType := d.Get("type").(string)
	globalReferenceId := d.Get("global_reference_id").(string)
	forTest := false
	isActive := true
	if val, ok := d.GetOk("for_test"); ok {
		forTest = val.(bool)
	}
	if val, ok := d.GetOk("is_active"); ok {
		isActive = val.(bool)
	}
	var deviceIdentifier iam.DeviceIdentifier
	if v, ok := d.GetOk("external_identifier"); ok {
		vL := v.(*schema.Set).List()
		var aVi []interface{}
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			deviceIdentifier.System = mVi["system"].(string)
			deviceIdentifier.Value = mVi["value"].(string)
			aVi = mVi["type"].(*schema.Set).List()
		}
		// Code readout
		for _, vi := range aVi {
			mVi := vi.(map[string]interface{})
			deviceIdentifier.Type = iam.CodeableConcept{
				Code: mVi["code"].(string),
				Text: mVi["text"].(string),
			}
		}
	}
	resource := iam.Device{
		LoginID:  loginId,
		Password: password,
		ForTest:  forTest,
		IsActive: isActive,
		Type:     deviceType,
		//Text:              text,
		DeviceExtID:       deviceIdentifier,
		GlobalReferenceID: globalReferenceId,
		ApplicationID:     applicationId,
		OrganizationID:    organizationId,
	}
	if val, ok := d.GetOk("debug_until"); ok {
		debugUntil, err := time.Parse(val.(string), time.RFC3339)
		if err == nil {
			resource.DebugUntil = &debugUntil
		}
	}
	return resource
}

func deviceToSchema(resource iam.Device, d *schema.ResourceData) {
	_ = d.Set("login_id", resource.LoginID)
	_ = d.Set("is_active", resource.IsActive)
	_ = d.Set("for_test", resource.ForTest)
	_ = d.Set("global_reference_id", resource.GlobalReferenceID)
	if resource.RegistrationDate != nil {
		_ = d.Set("registration_date", resource.RegistrationDate.Format(time.RFC3339))
	}
	if resource.DebugUntil != nil {
		_ = d.Set("debug_until", resource.DebugUntil.Format(time.RFC3339))
	}
	_ = d.Set("text", resource.Text)
	_ = d.Set("organization_id", resource.OrganizationID)
	_ = d.Set("application_id", resource.ApplicationID)
	_ = d.Set("type", resource.Type)

	// Build type
	tc := &schema.Set{F: schema.HashResource(deviceTypeSchema())}
	tcMap := make(map[string]interface{})
	tcMap["code"] = resource.DeviceExtID.Type.Code
	tcMap["text"] = resource.DeviceExtID.Type.Text
	tc.Add(tcMap)

	// Build external_identifier
	externalId := make(map[string]interface{})
	externalId["value"] = resource.DeviceExtID.Value
	externalId["system"] = resource.DeviceExtID.System
	externalId["type"] = tc

	ei := &schema.Set{F: schema.HashResource(externalIdentifierSchema())}
	ei.Add(externalId)

	_ = d.Set("external_identifier", ei)
}

func resourceIAMDeviceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	var ok bool
	var resp *iam.Response

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var err error

		ok, resp, err = client.Devices.DeleteDevice(iam.Device{ID: d.Id()})
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if !ok {
		return diag.FromErr(fmt.Errorf("error in DeleteDevice('%s'): [%v] %w", d.Id(), resp, err))
	}
	d.SetId("")
	return diags
}

func resourceIAMDeviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	device := schemaToDevice(d)
	device.ID = d.Id()

	var updatedDevice *iam.Device
	var resp *iam.Response

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var err error

		updatedDevice, resp, err = client.Devices.UpdateDevice(device)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp != nil && resp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(fmt.Errorf("error in UpdateDevice('%s'): %w", d.Id(), err))
	}
	deviceToSchema(*updatedDevice, d)
	return diags
}

func resourceIAMDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var device *iam.Device
	var resp *iam.Response

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var err error
		device, resp, err = client.Devices.GetDeviceByID(d.Id())
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp != nil && (resp.StatusCode() == http.StatusNotFound || resp.StatusCode() == http.StatusForbidden) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(fmt.Errorf("error in GetDeviceById('%s'): %w", d.Id(), err))
	}
	deviceToSchema(*device, d)

	return diags
}

func resourceIAMDeviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	device := schemaToDevice(d)

	if device.GlobalReferenceID == "" {
		result, err := uuid.GenerateUUID()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error generating uuid: %w", err))
		}
		device.GlobalReferenceID = result
	}
	var createdDevice *iam.Device
	var resp *iam.Response

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var err error
		createdDevice, resp, err = client.Devices.CreateDevice(device)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode() != http.StatusConflict {
			return diag.FromErr(err)
		}
		createdDevices, _, err := client.Devices.GetDevices(&iam.GetDevicesOptions{
			LoginID:        &device.LoginID,
			OrganizationID: &device.OrganizationID,
		})
		if err != nil || len(*createdDevices) == 0 {
			return diag.FromErr(fmt.Errorf("GetDevices after 409 (len=%d): %w", len(*createdDevices), err))
		}
		createdDevice = &(*createdDevices)[0]
		if createdDevice.OrganizationID != device.OrganizationID {
			return diag.FromErr(fmt.Errorf("existing application found but organization mismatch: '%s' != '%s'", createdDevice.OrganizationID, createdDevice.OrganizationID))
		}
		// We found a matching existing application, go with it
	}
	if createdDevice == nil {
		return diag.FromErr(fmt.Errorf("unexpected failure creating '%s': [%v] [%v]", device.LoginID, err, resp))
	}
	d.SetId(createdDevice.ID)

	return diags
}
