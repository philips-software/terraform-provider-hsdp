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

func ResourceConnectMDMStandardService() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 3,
		CreateContext: resourceConnectMDMStandardServiceCreate,
		ReadContext:   resourceConnectMDMStandardServiceRead,
		UpdateContext: resourceConnectMDMStandardServiceUpdate,
		DeleteContext: resourceConnectMDMStandardServiceDelete,

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
			"trusted": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				MinItems: 1,
				MaxItems: 1,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_url": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 5,
				Elem:     serviceURLSchema(),
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

func serviceURLSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sort_order": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"authentication_method_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func schemaToStandardService(d *schema.ResourceData) mdm.StandardService {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	orgIdentifier := d.Get("organization_identifier").(string)
	tags := tools.ExpandStringList(d.Get("tags").(*schema.Set).List())
	trusted := d.Get("trusted").(bool)

	var serviceURLS []mdm.ServiceURL

	if v, ok := d.GetOk("service_url"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			var serviceURL mdm.ServiceURL
			serviceURL.URL = mVi["url"].(string)
			serviceURL.SortOrder = mVi["sort_order"].(int)
			if id := mVi["authentication_method_id"].(string); id != "" {
				serviceURL.AuthenticationMethodID = &mdm.Reference{
					Reference: id,
				}
			}
			serviceURLS = append(serviceURLS, serviceURL)
		}
	}

	service := mdm.StandardService{
		Name:        name,
		Description: description,
		ServiceUrls: serviceURLS,
		Trusted:     trusted,
		Tags:        tags,
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

func standardServiceToSchema(service mdm.StandardService, d *schema.ResourceData) {
	s := &schema.Set{F: schema.HashResource(serviceURLSchema())}
	for _, serviceURL := range service.ServiceUrls {
		entry := make(map[string]interface{})
		entry["url"] = serviceURL.URL
		entry["sort_order"] = serviceURL.SortOrder
		if serviceURL.AuthenticationMethodID != nil {
			entry["authentication_method_id"] = serviceURL.AuthenticationMethodID.Reference
		}
		s.Add(entry)
	}
	_ = d.Set("service_url", s)
	_ = d.Set("description", service.Description)
	_ = d.Set("name", service.Name)
	_ = d.Set("tags", service.Tags)
	_ = d.Set("guid", service.ID)
	if service.OrganizationGuid != nil && service.OrganizationGuid.Value != "" {
		value := service.OrganizationGuid.Value
		if service.OrganizationGuid.System != "" {
			value = fmt.Sprintf("%s|%s", service.OrganizationGuid.System, service.OrganizationGuid.Value)
		}
		_ = d.Set("organization_identifier", value)
	}
}

func resourceConnectMDMStandardServiceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	service := schemaToStandardService(d)

	var created *mdm.StandardService
	var resp *mdm.Response
	err = tools.TryHTTPCall(func() (*http.Response, error) {
		var err error
		created, resp, err = client.StandardServices.CreateStandardService(service)
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
		return diag.FromErr(fmt.Errorf("failed to create standard service: %d", resp.StatusCode))
	}
	_ = d.Set("guid", created.ID)
	d.SetId(fmt.Sprintf("StandardService/%s", created.ID))
	return resourceConnectMDMStandardServiceRead(ctx, d, m)
}

func resourceConnectMDMStandardServiceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	service, resp, err := client.StandardServices.GetStandardServiceByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	standardServiceToSchema(*service, d)
	return diags
}

func resourceConnectMDMStandardServiceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)

	service := schemaToStandardService(d)
	service.ID = id

	_, _, err = client.StandardServices.Update(service)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		return diags
	}
	return resourceConnectMDMStandardServiceRead(ctx, d, m)
}

func resourceConnectMDMStandardServiceDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	service, _, err := client.StandardServices.GetStandardServiceByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.StandardServices.DeleteStandardService(*service)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
