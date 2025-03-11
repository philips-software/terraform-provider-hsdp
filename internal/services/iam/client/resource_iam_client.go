package client

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"

	"github.com/dip-software/go-dip-api/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var descriptions = map[string]string{
	"client": "A client identity describes a version of a registered application, its access scopes, and its credentials. Typically, the client identity represents a foreground (or user-facing) application. A client acts on other resources.",
}

func ResourceIAMClient() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    ResourceIAMClientV0().CoreConfigSchema().ImpliedType(),
				Upgrade: patchIAMClientV0,
				Version: 0,
			},
		},
		CreateContext: resourceIAMClientCreate,
		ReadContext:   resourceIAMClientRead,
		UpdateContext: resourceIAMClientUpdate,
		DeleteContext: resourceIAMClientDelete,
		Description:   descriptions["client"],
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the client.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the client. Either 'Public' or 'Confidential'.",
			},
			"client_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressCaseDiffs,
				Description:      "The client id",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "The password to use (8-16 chars, at least one capital, number, special char).",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of the client.",
			},
			"application_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The application ID to attach this client to.",
			},
			"global_reference_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Reference identifier defined by the provisioning user. This reference Identifier will be carried over to identify the provisioned resource across deployment instances.",
			},
			"redirection_uris": {
				Type:        schema.TypeSet,
				MaxItems:    100,
				MinItems:    0,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of valid RedirectionURIs for this client.",
			},
			"response_types": {
				Type:        schema.TypeSet,
				MaxItems:    100,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Examples of response types are 'code id_token', 'token id_token', etc.",
			},
			"scopes": {
				Type:        schema.TypeSet,
				MaxItems:    100,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of supported scopes for this client.",
			},
			"default_scopes": {
				Type:        schema.TypeSet,
				MaxItems:    100,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Default scopes. You do not have to specify these explicitly when requesting a token.",
			},
			"consent_implied": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag when enabled, the resource owner will not be asked for consent during authorization flows.",
			},
			"access_token_lifetime": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1800,
				Description: "Lifetime of the access token in seconds. If not specified, system default life time (1800 secs) will be considered.",
			},
			"refresh_token_lifetime": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2592000,
				Description: "Lifetime of the refresh token in seconds. If not specified, system default life time (2592000 secs) will be considered.",
			},
			"id_token_lifetime": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
				Description: "Lifetime of the jwt token generated in case openid scope is enabled for the client. If not specified, system default life time (3600 secs) will be considered.",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if the client is disabled e.g. because the managing Organization is disabled.",
			},
		},
	}
}

func resourceIAMClientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	token, err := client.Token()
	if err != nil {
		return diag.FromErr(err)
	}
	if token == "" {
		return diag.FromErr(config.ErrMissingIAMCredentials)
	}

	var cl iam.ApplicationClient
	cl.Description = d.Get("description").(string)
	cl.Name = d.Get("name").(string)
	cl.ClientID = d.Get("client_id").(string)
	cl.Type = d.Get("type").(string)
	cl.GlobalReferenceID = d.Get("global_reference_id").(string)
	cl.Password = d.Get("password").(string)
	cl.Name = d.Get("name").(string)
	cl.RedirectionURIs = tools.ExpandStringList(d.Get("redirection_uris").(*schema.Set).List())
	cl.ResponseTypes = tools.ExpandStringList(d.Get("response_types").(*schema.Set).List())
	cl.ApplicationID = d.Get("application_id").(string)
	cl.Scopes = tools.ExpandStringList(d.Get("scopes").(*schema.Set).List())
	cl.DefaultScopes = tools.ExpandStringList(d.Get("default_scopes").(*schema.Set).List())
	cl.IDTokenLifetime = d.Get("id_token_lifetime").(int)
	cl.RefreshTokenLifetime = d.Get("refresh_token_lifetime").(int)
	cl.AccessTokenLifetime = d.Get("access_token_lifetime").(int)
	cl.ConsentImplied = d.Get("consent_implied").(bool)

	var createdClient *iam.ApplicationClient

	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		var resp *iam.Response
		createdClient, _, err = client.Clients.CreateClient(cl)
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
	d.SetId(createdClient.ID)
	_ = d.Set("password", cl.Password)
	return resourceIAMClientRead(ctx, d, m)
}

func resourceIAMClientRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	token, err := client.Token()
	if err != nil {
		return diag.FromErr(err)
	}
	if token == "" {
		return diag.FromErr(config.ErrMissingIAMCredentials)
	}

	id := d.Id()
	cl, resp, err := client.Clients.GetClientByID(id)
	if err != nil {
		if resp != nil && (resp.StatusCode() == http.StatusNotFound || resp.StatusCode() == http.StatusForbidden) {
			// TODO: we should retrieve the managing organization of the Proposition and check for `CLIENT.READ` permissions
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("description", cl.Description)
	_ = d.Set("name", cl.Name)
	_ = d.Set("client_id", cl.ClientID)
	_ = d.Set("type", cl.Type)
	_ = d.Set("application_id", cl.ApplicationID)
	_ = d.Set("global_reference_id", cl.GlobalReferenceID)
	_ = d.Set("redirection_uris", cl.RedirectionURIs)
	_ = d.Set("response_types", cl.ResponseTypes)
	_ = d.Set("scopes", cl.Scopes)
	_ = d.Set("default_scopes", cl.DefaultScopes)
	_ = d.Set("access_token_lifetime", cl.AccessTokenLifetime)
	_ = d.Set("refresh_token_lifetime", cl.RefreshTokenLifetime)
	_ = d.Set("id_token_lifetime", cl.IDTokenLifetime)
	_ = d.Set("disabled", cl.Disabled)
	_ = d.Set("consent_implied", cl.ConsentImplied)

	return diags
}

func resourceIAMClientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	token, err := client.Token()
	if err != nil {
		return diag.FromErr(err)
	}
	if token == "" {
		return diag.FromErr(config.ErrMissingIAMCredentials)
	}

	var cl iam.ApplicationClient
	cl.ID = d.Id()

	if d.HasChange("scopes") || d.HasChange("default_scopes") {
		newScopes := tools.ExpandStringList(d.Get("scopes").(*schema.Set).List())
		newDefaultScopes := tools.ExpandStringList(d.Get("default_scopes").(*schema.Set).List())
		if d.HasChange("scopes") {
			_, ns := d.GetChange("scopes")
			newScopes = tools.ExpandStringList(ns.(*schema.Set).List())
		}
		if d.HasChange("default_scopes") {
			_, nd := d.GetChange("default_scopes")
			newDefaultScopes = tools.ExpandStringList(nd.(*schema.Set).List())
		}
		_, _, err := client.Clients.UpdateScopes(cl, newScopes, newDefaultScopes)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("access_token_lifetime") ||
		d.HasChange("refresh_token_lifetime") ||
		d.HasChange("id_token_lifetime") ||
		d.HasChange("consent_implied") ||
		d.HasChange("global_reference_id") ||
		d.HasChange("response_types") ||
		d.HasChange("redirection_uris") ||
		d.HasChange("description") {
		cl, _, err := client.Clients.GetClientByID(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		cl.RedirectionURIs = tools.ExpandStringList(d.Get("redirection_uris").(*schema.Set).List())
		cl.ResponseTypes = tools.ExpandStringList(d.Get("response_types").(*schema.Set).List())
		cl.ConsentImplied = d.Get("consent_implied").(bool)
		cl.AccessTokenLifetime = d.Get("access_token_lifetime").(int)
		cl.RefreshTokenLifetime = d.Get("refresh_token_lifetime").(int)
		cl.IDTokenLifetime = d.Get("id_token_lifetime").(int)
		cl.GlobalReferenceID = d.Get("global_reference_id").(string)
		cl.Description = d.Get("description").(string)
		_, _, err = client.Clients.UpdateClient(*cl)
		if err != nil {
			return diag.FromErr(err)
		}
		return resourceIAMClientRead(ctx, d, m)
	}
	return diags
}

func resourceIAMClientDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	token, err := client.Token()
	if err != nil {
		return diag.FromErr(err)
	}
	if token == "" {
		return diag.FromErr(config.ErrMissingIAMCredentials)
	}

	var cl iam.ApplicationClient
	cl.ID = d.Id()
	var ok bool
	var resp *iam.Response

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		ok, resp, err = client.Clients.DeleteClient(cl)
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp != nil && resp.StatusCode() == http.StatusNotFound { // Gone already
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrDeleteClientFailed)
	}
	d.SetId("")
	return diags
}
