package mdm

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceConnectMDMDeviceGroup() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceConnectMDMDeviceGroupCreate,
		ReadContext:   resourceConnectMDMDeviceGroupRead,
		UpdateContext: resourceConnectMDMDeviceGroupUpdate,
		DeleteContext: resourceConnectMDMDeviceGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_iam_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: tools.SuppressMulti(
					tools.SuppressWhenGenerated,
					tools.SuppressDefaultSystemValue),
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

func schemaToDeviceGroup(d *schema.ResourceData) mdm.DeviceGroup {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	applicationId := d.Get("application_id").(string)
	defaultIAMGroupId := d.Get("default_iam_group_id").(string)

	resource := mdm.DeviceGroup{
		Name:          name,
		Description:   description,
		ApplicationId: mdm.Reference{Reference: applicationId},
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
	return resource
}

func deviceGroupToSchema(resource mdm.DeviceGroup, d *schema.ResourceData) {
	_ = d.Set("description", resource.Description)
	_ = d.Set("name", resource.Name)
	_ = d.Set("application_id", resource.ApplicationId.Reference)
	_ = d.Set("guid", resource.ID)
	if resource.DefaultGroupGuid != nil && resource.DefaultGroupGuid.Value != "" {
		value := resource.DefaultGroupGuid.Value
		if resource.DefaultGroupGuid.System != "" {
			value = fmt.Sprintf("%s|%s", resource.DefaultGroupGuid.System, resource.DefaultGroupGuid.Value)
		}
		_ = d.Set("default_iam_group_id", value)
	} else {
		_ = d.Set("default_iam_group_id", nil)
	}
}

func resourceConnectMDMDeviceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToDeviceGroup(d)

	var created *mdm.DeviceGroup
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.DeviceGroups.Create(resource)
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
	d.SetId(fmt.Sprintf("DeviceGroup/%s", created.ID))
	return resourceConnectMDMDeviceGroupRead(ctx, d, m)
}

func resourceConnectMDMDeviceGroupRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "DeviceGroup/%s", &id)
	resource, resp, err := client.DeviceGroups.GetByID(id)
	if err != nil {
		if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	deviceGroupToSchema(*resource, d)
	return diags
}

func resourceConnectMDMDeviceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)

	resource := schemaToDeviceGroup(d)
	resource.ID = id

	_, _, err = client.DeviceGroups.Update(resource)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		return diags
	}
	return resourceConnectMDMDeviceGroupRead(ctx, d, m)
}

func resourceConnectMDMDeviceGroupDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	resource, _, err := client.DeviceGroups.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.DeviceGroups.Delete(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
