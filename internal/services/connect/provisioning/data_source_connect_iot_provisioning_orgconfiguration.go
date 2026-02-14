package provisioning

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-dip-api/connect/provisioning"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func DataSourceConnectIoTProvisioningOrgConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectIoTProvisioningOrgConfigurationRead,
		Schema: map[string]*schema.Schema{
			"organization_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID of the organization",
			},
			"service_account": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Service account configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service account ID",
						},
						"service_account_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Service account private key",
						},
					},
				},
			},
			"bootstrap_signature": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Bootstrap signature configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Signature algorithm (e.g., RSA-SHA256)",
						},
						"public_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public key for bootstrap signature",
						},
						"config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Bootstrap signature config",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Signature type (RSA, ECC, DSA)",
									},
									"padding": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Padding type (e.g., RSA_PKCS1_PSS_PADDING)",
									},
									"salt_length": {
										Type:        schema.TypeString,
										Computed:    true,
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

func dataSourceConnectIoTProvisioningOrgConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	client, err := c.ProvisioningClient()
	if err != nil {
		return diag.FromErr(err)
	}

	organizationGuid := d.Get("organization_guid").(string)

	var resources *[]provisioning.OrgConfiguration
	var resp *provisioning.Response

	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resources, resp, err = client.OrgConfigurationsService.FindOrgConfiguration(&provisioning.GetOrgConfiguration{
			OrganizationGuid: &organizationGuid,
		})
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error finding Organization Configuration: %w", err))
	}

	if resources == nil || len(*resources) == 0 {
		d.SetId("")
		return diags
	}

	// Use the first result
	resource := (*resources)[0]

	d.SetId(resource.ID)
	orgConfigurationToSchema(resource, d)

	return diags
}
