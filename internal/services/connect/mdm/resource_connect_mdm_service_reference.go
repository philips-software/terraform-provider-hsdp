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

func ResourceConnectMDMServiceReference() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceConnectMDMServiceReferenceCreate,
		ReadContext:   resourceConnectMDMServiceReferenceRead,
		UpdateContext: resourceConnectMDMServiceReferenceUpdate,
		DeleteContext: resourceConnectMDMServiceReferenceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"standard_service_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"matching_rule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_action_ids": {
				Type:     schema.TypeSet,
				MaxItems: 20,
				MinItems: 1,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bootstrap_enabled": {
				Type:     schema.TypeBool,
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

func schemaToServiceReference(d *schema.ResourceData) mdm.ServiceReference {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	applicationId := d.Get("application_id").(string)
	standardServiceId := d.Get("standard_service_id").(string)
	matchingRule := d.Get("matching_rule").(string)
	serviceActionIDs := tools.ExpandStringList(d.Get("service_action_ids").(*schema.Set).List())
	bootstrapEnabled := d.Get("bootstrap_enabled").(bool)

	resource := mdm.ServiceReference{
		Name:              name,
		Description:       description,
		ApplicationID:     mdm.Reference{Reference: applicationId},
		StandardServiceID: mdm.Reference{Reference: standardServiceId},
		MatchingRule:      matchingRule,
		BootstrapEnabled:  bootstrapEnabled,
	}
	for _, action := range serviceActionIDs {
		resource.ServiceActionIDs = append(resource.ServiceActionIDs, mdm.Reference{Reference: action})
	}
	return resource
}

func ServiceReferenceToSchema(resource mdm.ServiceReference, d *schema.ResourceData) {
	_ = d.Set("description", resource.Description)
	_ = d.Set("name", resource.Name)
	_ = d.Set("application_id", resource.ApplicationID.Reference)
	_ = d.Set("standard_service_id", resource.StandardServiceID.Reference)
	_ = d.Set("matching_rule", resource.MatchingRule)
	_ = d.Set("guid", resource.ID)
	_ = d.Set("bootstrap_enabled", resource.BootstrapEnabled)
	var serviceActionIDs []string
	for _, a := range resource.ServiceActionIDs {
		serviceActionIDs = append(serviceActionIDs, a.Reference)
	}
	_ = d.Set("service_action_ids", serviceActionIDs)
	_ = d.Set("guid", resource.ID)
}

func resourceConnectMDMServiceReferenceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToServiceReference(d)

	var created *mdm.ServiceReference
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.ServiceReferences.Create(resource)
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
	d.SetId(fmt.Sprintf("ServiceReference/%s", created.ID))
	return resourceConnectMDMServiceReferenceRead(ctx, d, m)
}

func resourceConnectMDMServiceReferenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "ServiceReference/%s", &id)
	var resource *mdm.ServiceReference
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.ServiceReferences.GetByID(id)
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
	ServiceReferenceToSchema(*resource, d)
	return diags
}

func resourceConnectMDMServiceReferenceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)

	resource := schemaToServiceReference(d)
	resource.ID = id

	_, _, err = client.ServiceReferences.Update(resource)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		return diags
	}
	return resourceConnectMDMServiceReferenceRead(ctx, d, m)
}

func resourceConnectMDMServiceReferenceDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	resource, _, err := client.ServiceReferences.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.ServiceReferences.Delete(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
