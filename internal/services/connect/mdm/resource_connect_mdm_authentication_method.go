package mdm

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-dip-api/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceConnectMDMAuthenticationMethod() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceConnectMDMAuthenticationMethodCreate,
		ReadContext:   resourceConnectMDMAuthenticationMethodRead,
		UpdateContext: resourceConnectMDMAuthenticationMethodUpdate,
		DeleteContext: resourceConnectMDMAuthenticationMethodDelete,

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
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},
			"auth_method": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organization_id": {
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

func schemaToAuthenticationMethod(d *schema.ResourceData) mdm.AuthenticationMethod {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	apiVersion := d.Get("api_version").(string)
	organizationId := d.Get("organization_id").(string)
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	loginName := d.Get("login_name").(string)
	password := d.Get("password").(string)
	authURL := d.Get("auth_url").(string)
	authMethod := d.Get("auth_method").(string)

	resource := mdm.AuthenticationMethod{
		Name:         name,
		Description:  description,
		APIVersion:   apiVersion,
		LoginName:    loginName,
		Password:     password,
		ClientID:     clientId,
		ClientSecret: clientSecret,
		AuthMethod:   authMethod,
		AuthURL:      authURL,
	}
	if len(organizationId) > 0 {
		identifier := mdm.Identifier{}
		parts := strings.Split(organizationId, "|")
		if len(parts) > 1 {
			identifier.System = parts[0]
			identifier.Value = parts[1]
		} else {
			identifier.Value = organizationId
		}
		resource.OrganizationGuid = &identifier
	}
	return resource
}

func AuthenticationMethodToSchema(resource mdm.AuthenticationMethod, d *schema.ResourceData) {
	_ = d.Set("description", resource.Description)
	_ = d.Set("name", resource.Name)
	_ = d.Set("login_name", resource.LoginName)
	_ = d.Set("password", resource.Password)
	_ = d.Set("client_id", resource.ClientID)
	_ = d.Set("client_secret", resource.ClientSecret)
	_ = d.Set("auth_method", resource.AuthMethod)
	_ = d.Set("auth_url", resource.AuthURL)
	_ = d.Set("api_version", resource.APIVersion)

	_ = d.Set("guid", resource.ID)
	if resource.OrganizationGuid != nil && resource.OrganizationGuid.Value != "" {
		value := resource.OrganizationGuid.Value
		if resource.OrganizationGuid.System != "" {
			value = fmt.Sprintf("%s|%s", resource.OrganizationGuid.System, resource.OrganizationGuid.Value)
		}
		_ = d.Set("organization_id", value)
	} else {
		_ = d.Set("organization_id", nil)
	}
}

func resourceConnectMDMAuthenticationMethodCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToAuthenticationMethod(d)

	var created *mdm.AuthenticationMethod
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.AuthenticationMethods.Create(resource)
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
	d.SetId(fmt.Sprintf("AuthenticationMethod/%s", created.ID))
	return resourceConnectMDMAuthenticationMethodRead(ctx, d, m)
}

func resourceConnectMDMAuthenticationMethodRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "AuthenticationMethod/%s", &id)
	var resource *mdm.AuthenticationMethod
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.AuthenticationMethods.GetByID(id)
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
	AuthenticationMethodToSchema(*resource, d)
	return diags
}

func resourceConnectMDMAuthenticationMethodUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)

	resource := schemaToAuthenticationMethod(d)
	resource.ID = id

	_, _, err = client.AuthenticationMethods.Update(resource)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		return diags
	}
	return resourceConnectMDMAuthenticationMethodRead(ctx, d, m)
}

func resourceConnectMDMAuthenticationMethodDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	resource, _, err := client.AuthenticationMethods.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.AuthenticationMethods.Delete(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
