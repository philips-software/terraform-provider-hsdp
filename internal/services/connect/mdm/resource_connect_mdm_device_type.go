package mdm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dip-software/go-dip-api/connect/mdm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceConnectMDMDeviceType() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceConnectMDMDeviceTypeCreate,
		ReadContext:   resourceConnectMDMDeviceTypeRead,
		UpdateContext: resourceConnectMDMDeviceTypeUpdate,
		DeleteContext: resourceConnectMDMDeviceTypeDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"commercial_type_number": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"device_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"default_iam_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: tools.SuppressMulti(
					tools.SuppressWhenGenerated,
					tools.SuppressDefaultSystemValue),
			},
			"custom_type_attributes": {
				Type:     schema.TypeMap,
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

func schemaToDeviceType(d *schema.ResourceData) mdm.DeviceType {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	deviceGroupId := d.Get("device_group_id").(string)
	defaultIAMGroupId := d.Get("default_iam_group_id").(string)
	ctn := d.Get("commercial_type_number").(string)

	resource := mdm.DeviceType{
		Name:          name,
		Description:   description,
		DeviceGroupId: mdm.Reference{Reference: deviceGroupId},
		CTN:           ctn,
	}
	if len(defaultIAMGroupId) > 0 {
		identifier := mdm.Identifier{}
		parts := strings.Split(defaultIAMGroupId, "|")
		if len(parts) > 1 {
			identifier.System = parts[0]
			identifier.Value = parts[1]
		} else {
			identifier.Value = defaultIAMGroupId
		}
		resource.DefaultGroupGuid = &identifier
	}
	if e, ok := d.GetOk("custom_type_attributes"); ok {
		custom := make(map[string]string)
		if env, ok := e.(map[string]interface{}); ok {
			for k, v := range env {
				custom[k] = v.(string)
			}
			raw, err := json.Marshal(custom)
			if err == nil {
				resource.CustomTypeAttributes = raw
			}
		}
	}
	return resource
}

func DeviceTypeToSchema(resource mdm.DeviceType, d *schema.ResourceData) {
	_ = d.Set("description", resource.Description)
	_ = d.Set("name", resource.Name)
	_ = d.Set("device_group_id", resource.DeviceGroupId.Reference)
	_ = d.Set("guid", resource.ID)
	_ = d.Set("commercial_type_number", resource.CTN)
	if resource.DefaultGroupGuid != nil && resource.DefaultGroupGuid.Value != "" {
		value := resource.DefaultGroupGuid.Value
		if resource.DefaultGroupGuid.System != "" {
			value = fmt.Sprintf("%s|%s", resource.DefaultGroupGuid.System, resource.DefaultGroupGuid.Value)
		}
		_ = d.Set("default_iam_group_id", value)
	}
	var custom map[string]interface{}
	if err := json.Unmarshal(resource.CustomTypeAttributes, &custom); err == nil {
		_ = d.Set("custom_type_attributes", custom)
	}
}

func resourceConnectMDMDeviceTypeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToDeviceType(d)

	var created *mdm.DeviceType
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.DeviceTypes.Create(resource)
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
	d.SetId(fmt.Sprintf("DeviceType/%s", created.ID))
	return resourceConnectMDMDeviceTypeRead(ctx, d, m)
}

func resourceConnectMDMDeviceTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "DeviceType/%s", &id)
	var resource *mdm.DeviceType
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.DeviceTypes.GetByID(id)
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
	DeviceTypeToSchema(*resource, d)
	return diags
}

func resourceConnectMDMDeviceTypeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)

	resource := schemaToDeviceType(d)
	resource.ID = id

	_, _, err = client.DeviceTypes.Update(resource)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		return diags
	}
	return resourceConnectMDMDeviceTypeRead(ctx, d, m)
}

func resourceConnectMDMDeviceTypeDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	resource, _, err := client.DeviceTypes.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.DeviceTypes.Delete(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
