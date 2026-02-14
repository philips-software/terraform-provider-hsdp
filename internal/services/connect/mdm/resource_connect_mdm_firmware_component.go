package mdm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/philips-software/go-dip-api/connect/mdm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceConnectMDMFirmwareComponent() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceConnectMDMFirmwareComponentCreate,
		ReadContext:   resourceConnectMDMFirmwareComponentRead,
		UpdateContext: resourceConnectMDMFirmwareComponentUpdate,
		DeleteContext: resourceConnectMDMFirmwareComponentDelete,

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
			"device_type_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"main_component": {
				Type:     schema.TypeBool,
				Default:  true,
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

func schemaToFirmwareComponent(d *schema.ResourceData) mdm.FirmwareComponent {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	deviceTypeId := d.Get("device_type_id").(string)
	mainComponent := d.Get("main_component").(bool)

	resource := mdm.FirmwareComponent{
		Name:          name,
		Description:   description,
		DeviceTypeId:  mdm.Reference{Reference: deviceTypeId},
		MainComponent: mainComponent,
	}
	return resource
}

func FirmwareComponentToSchema(resource mdm.FirmwareComponent, d *schema.ResourceData) {
	_ = d.Set("name", resource.Name)
	_ = d.Set("description", resource.Description)
	_ = d.Set("device_type_id", resource.DeviceTypeId.Reference)
	_ = d.Set("main_component", resource.MainComponent)
	_ = d.Set("guid", resource.ID)
}

func resourceConnectMDMFirmwareComponentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToFirmwareComponent(d)

	var created *mdm.FirmwareComponent
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.FirmwareComponents.Create(resource)
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
	d.SetId(fmt.Sprintf("FirmwareComponent/%s", created.ID))
	return resourceConnectMDMFirmwareComponentRead(ctx, d, m)
}

func resourceConnectMDMFirmwareComponentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "FirmwareComponent/%s", &id)
	var resource *mdm.FirmwareComponent
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.FirmwareComponents.GetByID(id)
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
	FirmwareComponentToSchema(*resource, d)
	return diags
}

func resourceConnectMDMFirmwareComponentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)

	service := schemaToFirmwareComponent(d)
	service.ID = id

	_, _, err = client.FirmwareComponents.Update(service)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		return diags
	}
	return resourceConnectMDMFirmwareComponentRead(ctx, d, m)
}

func resourceConnectMDMFirmwareComponentDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	resource, _, err := client.FirmwareComponents.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.FirmwareComponents.Delete(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
