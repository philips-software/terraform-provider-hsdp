package provisioning

import (
	"context"
	"fmt"
	"net/http"

	"github.com/philips-software/go-dip-api/connect/provisioning"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceConnectIoTProvisioningOrgConfiguration() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		CreateContext: resourceConnectIoTProvisioningOrgConfigurationCreate,
		ReadContext:   resourceConnectIoTProvisioningOrgConfigurationRead,
		UpdateContext: resourceConnectIoTProvisioningOrgConfigurationUpdate,
		DeleteContext: resourceConnectIoTProvisioningOrgConfigurationDelete,

		Schema: map[string]*schema.Schema{
			"organization_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The GUID of the organization",
			},
			"service_account": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Service account configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Service account ID",
						},
						"service_account_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Service account private key",
						},
					},
				},
			},
			"bootstrap_signature": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Bootstrap signature configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Signature algorithm (e.g., RSA-SHA256)",
						},
						"public_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Public key for bootstrap signature",
						},
						"config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Bootstrap signature config",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Signature type (RSA, ECC, DSA)",
									},
									"padding": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Padding type (e.g., RSA_PKCS1_PSS_PADDING)",
									},
									"salt_length": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Salt length (e.g., RSA_PSS_SALTLEN_MAX_SIGN)",
									},
								},
							},
						},
					},
				},
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the organization configuration",
			},
		},
	}
}

func schemaToOrgConfiguration(d *schema.ResourceData) provisioning.OrgConfiguration {
	organizationGuid := d.Get("organization_guid").(string)

	// Parse service account
	serviceAccountList := d.Get("service_account").([]interface{})
	var serviceAccount provisioning.ServiceAccount
	if len(serviceAccountList) > 0 {
		sa := serviceAccountList[0].(map[string]interface{})
		serviceAccount = provisioning.ServiceAccount{
			ServiceAccountId:  sa["service_account_id"].(string),
			ServiceAccountKey: sa["service_account_key"].(string),
		}
	}

	// Parse bootstrap signature
	bootstrapSignatureList := d.Get("bootstrap_signature").([]interface{})
	var bootstrapSignature provisioning.BootstrapSignature
	if len(bootstrapSignatureList) > 0 {
		bs := bootstrapSignatureList[0].(map[string]interface{})
		bootstrapSignature = provisioning.BootstrapSignature{
			Algorithm: bs["algorithm"].(string),
			PublicKey: bs["public_key"].(string),
		}

		// Parse config if present
		configList := bs["config"].([]interface{})
		if len(configList) > 0 {
			config := configList[0].(map[string]interface{})
			bootstrapSignature.Config = provisioning.BootStrapSignatureConfig{
				Type:       config["type"].(string),
				Padding:    config["padding"].(string),
				SaltLength: config["salt_length"].(string),
			}
		}
	}

	return provisioning.OrgConfiguration{
		ResourceType:       "OrgConfiguration",
		OrganizationGuid:   organizationGuid,
		ServiceAccount:     serviceAccount,
		BootstrapSignature: bootstrapSignature,
	}
}

func orgConfigurationToSchema(resource provisioning.OrgConfiguration, d *schema.ResourceData) {
	_ = d.Set("organization_guid", resource.OrganizationGuid)
	_ = d.Set("id", resource.ID)

	// Set service account
	serviceAccount := []map[string]interface{}{
		{
			"service_account_id":  resource.ServiceAccount.ServiceAccountId,
			"service_account_key": resource.ServiceAccount.ServiceAccountKey,
		},
	}
	_ = d.Set("service_account", serviceAccount)

	// Set bootstrap signature
	bootstrapSignature := []map[string]interface{}{
		{
			"algorithm":  resource.BootstrapSignature.Algorithm,
			"public_key": resource.BootstrapSignature.PublicKey,
		},
	}

	// Add config if present
	if resource.BootstrapSignature.Config.Type != "" ||
		resource.BootstrapSignature.Config.Padding != "" ||
		resource.BootstrapSignature.Config.SaltLength != "" {
		config := []map[string]interface{}{
			{
				"type":        resource.BootstrapSignature.Config.Type,
				"padding":     resource.BootstrapSignature.Config.Padding,
				"salt_length": resource.BootstrapSignature.Config.SaltLength,
			},
		}
		bootstrapSignature[0]["config"] = config
	}

	_ = d.Set("bootstrap_signature", bootstrapSignature)
}

func resourceConnectIoTProvisioningOrgConfigurationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.ProvisioningClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToOrgConfiguration(d)

	var created *provisioning.OrgConfiguration
	var resp *provisioning.Response

	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.OrgConfigurationsService.CreateOrganizationConfiguration(resource)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		return diag.FromErr(fmt.Errorf("error creating Organization Configuration (%d): %w", resp.StatusCode(), err))
	}

	if created == nil {
		return diag.FromErr(fmt.Errorf("unexpected error creating Organization Configuration: %v", resp))
	}

	d.SetId(created.ID)
	return resourceConnectIoTProvisioningOrgConfigurationRead(ctx, d, m)
}

func resourceConnectIoTProvisioningOrgConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	client, err := c.ProvisioningClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	var resource *provisioning.OrgConfiguration
	var resp *provisioning.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.OrgConfigurationsService.GetOrganizationConfigurationByID(id)
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
			return diags
		}
		return diag.FromErr(err)
	}
	orgConfigurationToSchema(*resource, d)
	return diags
}

func resourceConnectIoTProvisioningOrgConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	client, err := c.ProvisioningClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToOrgConfiguration(d)
	resource.ID = d.Id()

	var updated *provisioning.OrgConfiguration
	var resp *provisioning.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		updated, resp, err = client.OrgConfigurationsService.UpdateOrganizationConfiguration(resource)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating Organization Configuration: %w", err))
	}

	orgConfigurationToSchema(*updated, d)
	return nil
}

func resourceConnectIoTProvisioningOrgConfigurationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	client, err := c.ProvisioningClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToOrgConfiguration(d)
	resource.ID = d.Id()

	var resp *provisioning.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		_, resp, err = client.OrgConfigurationsService.DeleteOrganizationConfiguration(resource)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp != nil && resp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(fmt.Errorf("error deleting Organization Configuration: %w", err))
	}

	d.SetId("")
	return diags
}
