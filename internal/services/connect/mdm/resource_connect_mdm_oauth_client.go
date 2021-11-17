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

func ResourceConnectMDMOAuthClient() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceConnectMDMOAuthClientCreate,
		ReadContext:   resourceConnectMDMOAuthClientRead,
		UpdateContext: resourceConnectMDMOAuthClientUpdate,
		DeleteContext: resourceConnectMDMOAuthClientDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"global_reference_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"redirection_uris": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				ForceNew: true,
				Elem:     tools.StringSchema(),
			},
			"response_types": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				ForceNew: true,
				Elem:     tools.StringSchema(),
			},
			"scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     tools.StringSchema(),
			},
			"default_scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     tools.StringSchema(),
			},
			"user_client": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"client_revoked": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"client_guid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: tools.SuppressMulti(
					tools.SuppressWhenGenerated,
					tools.SuppressDefaultSystemValue),
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"bootstrap_client_guid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: tools.SuppressMulti(
					tools.SuppressWhenGenerated,
					tools.SuppressDefaultSystemValue),
			},
			"bootstrap_client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bootstrap_client_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
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

func setScopes(client *mdm.Client, d *schema.ResourceData) error {

	scopes := tools.ExpandStringList(d.Get("scopes").(*schema.Set).List())
	defaultScopes := tools.ExpandStringList(d.Get("default_scopes").(*schema.Set).List())
	resourceId := d.Get("guid").(string)

	resource, _, err := client.OAuthClients.GetOAuthClientByID(resourceId)
	if err != nil {
		return fmt.Errorf("read OAuth client: %w", err)
	}
	scopeDictionary, _, err := client.OAuthClientScopes.GetOAuthClientScopes(nil)
	if err != nil {
		return fmt.Errorf("retrieving scope dictionary: %w", err)
	}
	var allowedScopes []string
	for _, s := range *scopeDictionary {
		allowedScopes = append(allowedScopes, s.Scope())
	}
	for _, scope := range scopes {
		if !tools.ContainsString(allowedScopes, scope) {
			return fmt.Errorf("scope '%s' not allowed", scope)
		}
	}
	_, _, err = client.OAuthClients.UpdateScopes(*resource, scopes, defaultScopes)
	if err != nil {
		return fmt.Errorf("updating scopes: %w", err)
	}
	return nil
}

func oAuthClientScopesToSchema(resource mdm.OAuthClient, d *schema.ResourceData) error {
	// TODO
	return nil
}

func schemaToOAuthClient(d *schema.ResourceData) mdm.OAuthClient {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	applicationId := d.Get("application_id").(string)
	globalReferenceId := d.Get("global_reference_id").(string)
	bootstrapClientGuid := d.Get("bootstrap_client_guid").(string)
	clientGuid := d.Get("client_guid").(string)
	clientRevoked := d.Get("client_revoked").(bool)
	userClient := d.Get("user_client").(bool)
	redirectionURIs := tools.ExpandStringList(d.Get("redirection_uris").(*schema.Set).List())
	responseTypes := tools.ExpandStringList(d.Get("response_types").(*schema.Set).List())

	resource := mdm.OAuthClient{
		Name:              name,
		Description:       description,
		ApplicationId:     mdm.Reference{Reference: applicationId},
		GlobalReferenceID: globalReferenceId,
		UserClient:        userClient,
		ClientRevoked:     clientRevoked,
		ResponseTypes:     responseTypes,
		RedirectionURIs:   redirectionURIs,
	}
	if len(bootstrapClientGuid) > 0 {
		identifier := mdm.Identifier{}
		parts := strings.Split(bootstrapClientGuid, "|")
		if len(parts) > 1 {
			identifier.System = parts[0]
			identifier.Value = parts[1]
		} else {
			identifier.Value = bootstrapClientGuid
		}
		resource.BootstrapClientGuid = &identifier
	}
	if len(clientGuid) > 0 {
		identifier := mdm.Identifier{}
		parts := strings.Split(clientGuid, "|")
		if len(parts) > 1 {
			identifier.System = parts[0]
			identifier.Value = parts[1]
		} else {
			identifier.Value = clientGuid
		}
		resource.ClientGuid = &identifier
	}
	return resource
}

func oAuthClientToSchema(resource mdm.OAuthClient, d *schema.ResourceData) {
	_ = d.Set("description", resource.Description)
	_ = d.Set("name", resource.Name)
	_ = d.Set("application_id", resource.ApplicationId.Reference)
	_ = d.Set("global_reference_id", resource.GlobalReferenceID)
	_ = d.Set("guid", resource.ID)
	if resource.BootstrapClientGuid != nil && resource.BootstrapClientGuid.Value != "" {
		value := resource.BootstrapClientGuid.Value
		if resource.BootstrapClientGuid.System != "" {
			value = fmt.Sprintf("%s|%s", resource.BootstrapClientGuid.System, resource.BootstrapClientGuid.Value)
		}
		_ = d.Set("bootstrap_client_guid", value)
	}
	if resource.ClientGuid != nil && resource.ClientGuid.Value != "" {
		value := resource.ClientGuid.Value
		if resource.ClientGuid.System != "" {
			value = fmt.Sprintf("%s|%s", resource.ClientGuid.System, resource.ClientGuid.Value)
		}
		_ = d.Set("client_guid", value)
	}
	_ = d.Set("client_revoked", resource.ClientRevoked)
	_ = d.Set("user_client", resource.UserClient)
	_ = d.Set("redirection_uris", resource.RedirectionURIs)
	_ = d.Set("response_types", resource.ResponseTypes)
	_ = d.Set("client_id", resource.ClientID)
	_ = d.Set("client_secret", resource.ClientSecret)
	_ = d.Set("bootstrap_client_id", resource.BootstrapClientID)
	_ = d.Set("bootstrap_client_secret", resource.BootstrapClientSecret)
}

func resourceConnectMDMOAuthClientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToOAuthClient(d)

	var created *mdm.OAuthClient
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 20, func() (*http.Response, error) {
		var err error
		created, resp, err = client.OAuthClients.CreateOAuthClient(resource)
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
	d.SetId(fmt.Sprintf("OAuthClient/%s", created.ID))

	if err := setScopes(client, d); err != nil {
		// Clean up
		//._, _, _ = client.OAuthClients.DeleteOAuthClient(*created)
		//d.SetId("")
		//return diag.FromErr(err)
		diags := resourceConnectMDMOAuthClientRead(ctx, d, m)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to set scopes",
			Detail:   err.Error(),
		})
		return diags
	}

	return resourceConnectMDMOAuthClientRead(ctx, d, m)
}

func resourceConnectMDMOAuthClientRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "OAuthClient/%s", &id)
	resource, resp, err := client.OAuthClients.GetOAuthClientByID(id)
	if err != nil {
		if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	oAuthClientToSchema(*resource, d)
	_ = oAuthClientScopesToSchema(*resource, d)
	return diags
}

func resourceConnectMDMOAuthClientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	if !(d.HasChange("scopes") || d.HasChange("default_scopes")) {
		return diag.FromErr(fmt.Errorf("only 'scopes' and 'default_scopes' can be updated. this is a bug"))
	}
	if err := setScopes(client, d); err != nil {
		return diag.FromErr(err)
	}
	return resourceConnectMDMOAuthClientRead(ctx, d, m)
}

func resourceConnectMDMOAuthClientDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	resource, _, err := client.OAuthClients.GetOAuthClientByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.OAuthClients.DeleteOAuthClient(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
