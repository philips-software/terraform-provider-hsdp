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

func ResourceConnectMDMServiceAction() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 3,
		CreateContext: resourceConnectMDMServiceActionCreate,
		ReadContext:   resourceConnectMDMServiceActionRead,
		UpdateContext: resourceConnectMDMServiceActionUpdate,
		DeleteContext: resourceConnectMDMServiceActionDelete,

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
			"standard_service_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"organization_identifier": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tools.SuppressWhenGenerated,
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

func schemaToServiceAction(d *schema.ResourceData) mdm.ServiceAction {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	orgIdentifier := d.Get("organization_identifier").(string)
	standardServiceId := d.Get("standard_service_id").(string)

	service := mdm.ServiceAction{
		Name:              name,
		Description:       description,
		StandardServiceId: mdm.Reference{Reference: standardServiceId},
	}
	if len(orgIdentifier) > 0 {
		orgGuid := mdm.Identifier{}
		parts := strings.Split(orgIdentifier, "|")
		if len(parts) > 1 {
			orgGuid.System = parts[0]
			orgGuid.Value = parts[1]
		} else {
			orgGuid.Value = orgIdentifier
		}
		service.OrganizationGuid = &orgGuid
	}
	return service
}

func serviceActionToSchema(service mdm.ServiceAction, d *schema.ResourceData) {
	_ = d.Set("description", service.Description)
	_ = d.Set("name", service.Name)
	_ = d.Set("standard_service_id", service.StandardServiceId.Reference)
	_ = d.Set("guid", service.ID)
	if service.OrganizationGuid != nil && service.OrganizationGuid.Value != "" {
		value := service.OrganizationGuid.Value
		if service.OrganizationGuid.System != "" {
			value = fmt.Sprintf("%s|%s", service.OrganizationGuid.System, service.OrganizationGuid.Value)
		}
		_ = d.Set("organization_identifier", value)
	}
}

func resourceConnectMDMServiceActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	service := schemaToServiceAction(d)

	var created *mdm.ServiceAction
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.ServiceActions.Create(service)
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
		return diag.FromErr(fmt.Errorf("failed to create a service action: %d", resp.StatusCode))
	}
	_ = d.Set("guid", created.ID)
	d.SetId(fmt.Sprintf("ServiceAction/%s", created.ID))
	return resourceConnectMDMServiceActionRead(ctx, d, m)
}

func resourceConnectMDMServiceActionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	service, resp, err := client.ServiceActions.GetByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	serviceActionToSchema(*service, d)
	return diags
}

func resourceConnectMDMServiceActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)

	service := schemaToServiceAction(d)
	service.ID = id

	_, _, err = client.ServiceActions.Update(service)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		return diags
	}
	return resourceConnectMDMServiceActionRead(ctx, d, m)
}

func resourceConnectMDMServiceActionDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	service, _, err := client.ServiceActions.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.ServiceActions.Delete(*service)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
