package mdm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceConnectMDMOAuthClient() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
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
			"bootstrap_client_scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Optional: true,
				Elem:     tools.StringSchema(),
			},
			"bootstrap_client_default_scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Optional: true,
				Elem:     tools.StringSchema(),
			},
			"iam_scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Optional: true,
				Elem:     tools.StringSchema(),
			},
			"iam_default_scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Optional: true,
				Elem:     tools.StringSchema(),
			},
			"bootstrap_client_iam_scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Optional: true,
				Elem:     tools.StringSchema(),
			},
			"bootstrap_client_iam_default_scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Optional: true,
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
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Will be removed",
			},
			"client_guid_system": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_guid_value": {
				Type:     schema.TypeString,
				Computed: true,
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
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Will be removed",
			},
			"bootstrap_client_guid_system": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bootstrap_client_guid_value": {
				Type:     schema.TypeString,
				Computed: true,
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

func setScopes(client *mdm.Client, iamClient *iam.Client, d *schema.ResourceData) error {

	scopes := tools.ExpandStringList(d.Get("scopes").(*schema.Set).List())
	defaultScopes := tools.ExpandStringList(d.Get("default_scopes").(*schema.Set).List())

	bootstrapScopes := tools.ExpandStringList(d.Get("bootstrap_client_scopes").(*schema.Set).List())
	bootstrapDefaultScopes := tools.ExpandStringList(d.Get("bootstrap_client_default_scopes").(*schema.Set).List())
	resourceId := d.Get("guid").(string)

	resource, _, err := client.OAuthClients.GetOAuthClientByID(resourceId)
	if err != nil {
		return fmt.Errorf("read OAuth client: %w", err)
	}
	if resource == nil || resource.ClientGuid == nil {
		return fmt.Errorf("setScopes: missing IAM client GUID")
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

	for _, scope := range bootstrapScopes {
		if !tools.ContainsString(allowedScopes, scope) {
			return fmt.Errorf("bootstrap scope '%s' not allowed", scope)
		}
	}

	_, _, err = client.OAuthClients.UpdateScopes(*resource, scopes, defaultScopes)
	if err != nil {
		return fmt.Errorf("updating scopes: %w", err)
	}

	_, _, err = client.OAuthClients.UpdateScopesByFlag(*resource, bootstrapScopes, bootstrapDefaultScopes, true)
	if err != nil {
		return fmt.Errorf("updating bootstrap client scopes: %w", err)
	}
	// Set IAM scopes
	iamScopes := tools.ExpandStringList(d.Get("iam_scopes").(*schema.Set).List())
	iamDefaultScopes := tools.ExpandStringList(d.Get("iam_default_scopes").(*schema.Set).List())
	iamResource, _, err := iamClient.Clients.GetClientByID(resource.ClientGuid.Value)
	if err != nil {
		return fmt.Errorf("get IAM OAuth client: %v", err)
	}
	// Merge MDM scopes and IAM scopes
	combinedDefaultScopes := append(defaultScopes, iamDefaultScopes...)
	combinedScopes := append(scopes, iamScopes...)
	_, _, err = iamClient.Clients.UpdateScopes(*iamResource, combinedScopes, combinedDefaultScopes)

	if err != nil {
		return fmt.Errorf("updating IAM OAuth client scopes: %v", err)
	}

	bootstrapIamScopes := tools.ExpandStringList(d.Get("bootstrap_client_iam_scopes").(*schema.Set).List())
	bootstrapIamDefaultScopes := tools.ExpandStringList(d.Get("bootstrap_client_iam_default_scopes").(*schema.Set).List())
	bootstrapIamResource, _, err := iamClient.Clients.GetClientByID(resource.BootstrapClientGuid.Value)
	if err != nil {
		return fmt.Errorf("get IAM OAuth bootstrap client: %v", err)
	}

	// Merge MDM scopes and IAM scopes
	combinedBootstrapDefaultScopes := append(bootstrapDefaultScopes, bootstrapIamDefaultScopes...)
	combinedBootstrapScopes := append(bootstrapScopes, bootstrapIamScopes...)
	_, _, err = iamClient.Clients.UpdateScopes(*bootstrapIamResource, combinedBootstrapScopes, combinedBootstrapDefaultScopes)

	return err
}

func oAuthClientScopesToSchema(iamClient *iam.Client, resource mdm.OAuthClient, d *schema.ResourceData) error {
	if resource.ClientGuid == nil {
		return fmt.Errorf("missing IAM client GUID")
	}
	iamResource, _, err := iamClient.Clients.GetClientByID(resource.ClientGuid.Value)
	if err != nil {
		return fmt.Errorf("error retrieving IAM client: %v", err)
	}
	mdmScopes := tools.ExpandStringList(d.Get("scopes").(*schema.Set).List())
	mdmDefaultScopes := tools.ExpandStringList(d.Get("default_scopes").(*schema.Set).List())

	prunedIAMScopes := tools.Difference(iamResource.Scopes, mdmScopes)
	prunedIAMDefaultScopes := tools.Difference(iamResource.DefaultScopes, mdmDefaultScopes)

	_ = d.Set("iam_scopes", prunedIAMScopes)
	_ = d.Set("iam_default_scopes", prunedIAMDefaultScopes)
	_ = d.Set("scopes", tools.Difference(mdmScopes, prunedIAMScopes))
	_ = d.Set("default_scopes", tools.Difference(mdmDefaultScopes, prunedIAMDefaultScopes))
	return nil
}

func schemaToOAuthClient(d *schema.ResourceData) mdm.OAuthClient {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	applicationId := d.Get("application_id").(string)
	globalReferenceId := d.Get("global_reference_id").(string)
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
	return resource
}

func oAuthClientToSchema(resource mdm.OAuthClient, d *schema.ResourceData) {
	_ = d.Set("description", resource.Description)
	_ = d.Set("name", resource.Name)
	_ = d.Set("application_id", resource.ApplicationId.Reference)
	_ = d.Set("global_reference_id", resource.GlobalReferenceID)
	_ = d.Set("guid", resource.ID)
	if resource.BootstrapClientGuid != nil && resource.BootstrapClientGuid.Value != "" {
		_ = d.Set("bootstrap_client_guid_system", resource.BootstrapClientGuid.System)
		_ = d.Set("bootstrap_client_guid_value", resource.BootstrapClientGuid.Value)
	}
	if resource.ClientGuid != nil && resource.ClientGuid.Value != "" {
		_ = d.Set("client_guid_system", resource.ClientGuid.System)
		_ = d.Set("client_guid_value", resource.ClientGuid.Value)
	}
	_ = d.Set("client_revoked", resource.ClientRevoked)
	_ = d.Set("user_client", resource.UserClient)
	_ = d.Set("redirection_uris", resource.RedirectionURIs)
	_ = d.Set("response_types", resource.ResponseTypes)
	_ = d.Set("client_id", resource.ClientID)
	_ = d.Set("bootstrap_client_id", resource.BootstrapClientID)
	if resource.BootstrapClientSecret != "" {
		_ = d.Set("bootstrap_client_secret", resource.BootstrapClientSecret)
	}
	if resource.ClientSecret != "" {
		_ = d.Set("client_secret", resource.ClientSecret)
	}
}

func resourceConnectMDMOAuthClientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	iamClient, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	user_client := d.Get("user_client").(bool)
	bootstrap_client_scopes := tools.ExpandStringList(d.Get("bootstrap_client_scopes").(*schema.Set).List())
	bootstrap_client_default_scopes := tools.ExpandStringList(d.Get("bootstrap_client_default_scopes").(*schema.Set).List())
	bootstrap_client_iam_scopes := tools.ExpandStringList(d.Get("bootstrap_client_iam_scopes").(*schema.Set).List())
	bootstrap_client_iam_default_scopes := tools.ExpandStringList(d.Get("bootstrap_client_iam_default_scopes").(*schema.Set).List())

	if user_client && (len(bootstrap_client_scopes) > 0 || len(bootstrap_client_default_scopes) > 0 ||
		len(bootstrap_client_iam_scopes) > 0 || len(bootstrap_client_iam_default_scopes) > 0) {
		return diag.FromErr(fmt.Errorf("bootstrap client scopes are only allowed for non user clients"))
	}

	resource := schemaToOAuthClient(d)

	var created *mdm.OAuthClient
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
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
	d.SetId(fmt.Sprintf("OAuthClient/%s", created.ID))
	_ = d.Set("guid", created.ID)
	_ = d.Set("bootstrap_client_secret", created.BootstrapClientSecret)
	_ = d.Set("client_secret", created.ClientSecret)

	if err := setScopes(client, iamClient, d); err != nil {
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

func resourceConnectMDMOAuthClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	iamClient, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "OAuthClient/%s", &id)
	var resource *mdm.OAuthClient
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.OAuthClients.GetOAuthClientByID(id)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	oAuthClientToSchema(*resource, d)
	if err := oAuthClientScopesToSchema(iamClient, *resource, d); err != nil {
		return diag.FromErr(fmt.Errorf("oAuthClientScopesToSchema: %v", err))
	}
	return diags
}

func resourceConnectMDMOAuthClientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	iamClient, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	if !(d.HasChange("scopes") || d.HasChange("default_scopes") || d.HasChange("iam_scopes") || d.HasChange("iam_default_scopes") ||
		d.HasChange("bootstrap_client_scopes") || d.HasChange("bootstrap_client_default_scopes") ||
		d.HasChange("bootstrap_client_iam_scopes") || d.HasChange("bootstrap_client_iam_default_scopes")) {
		return diag.FromErr(fmt.Errorf("only 'scopes', 'default_scopes', 'iam_scopes', 'iam_default_scopes', " +
			"'bootstrap_client_scopes', 'bootstrap_client_default_scopes', 'bootstrap_client_iam_scopes' " +
			"or 'bootstrap_client_iam_default_scopes' can be updated. this is a bug"))
	}
	if err := setScopes(client, iamClient, d); err != nil {
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
