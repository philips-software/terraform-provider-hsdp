package hsdp

import (
	"context"
	"encoding/json"
	"os"

	"github.com/philips-software/terraform-provider-hsdp/internal/services/connect/dbs"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/connect/provisioning"

	"github.com/philips-software/terraform-provider-hsdp/internal/services/blr"

	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam/group_membership"

	"github.com/google/fhir/go/fhirversion"
	"github.com/google/fhir/go/jsonformat"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/configuration"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/discovery"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/edge"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam/application"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam/client"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam/device"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam/email_template"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam/group"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam/organization"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam/proposition"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam/role"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam/role_sharing_policy"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam/service"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam/user"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/tenant"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

const (
	Region           = "HSDP_REGION"
	Environment      = "HSDP_ENVIRONMENT"
	ServiceID        = "HSDP_IAM_SERVICE_ID"
	ServicePK        = "HSDP_IAM_SERVICE_PRIVATE_KEY"
	OrgAdminUsername = "HSDP_IAM_ORG_ADMIN_USERNAME"
	OrgAdminPassword = "HSDP_IAM_ORG_ADMIN_PASSWORD"
	ClientID         = "HSDP_IAM_OAUTH2_CLIENT_ID"
	ClientPassword   = "HSDP_IAM_OAUTH2_PASSWORD"
	SharedKey        = "HSDP_SHARED_KEY"
	SecretKey        = "HSDP_SECRET_KEY"
	UAAUsername      = "HSDP_UAA_USERNAME"
	UAAPassword      = "HSDP_UAA_PASSWORD"
	DebugLog         = "HSDP_DEBUG_LOG"
	DebugStdErr      = "HSDP_DEBUG_STDERR"
)

// Provider returns an instance of the HSDP provider
func Provider(build string) *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(Region, "us-east"),
				Description:  descriptions["region"],
				ValidateFunc: tools.ValidateRegion,
			},
			"environment": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(Environment, "client-test"),
				Description:  descriptions["environment"],
				ValidateFunc: tools.ValidateEnvironment,
			},
			"iam_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["iam_url"],
			},
			"idm_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["idm_url"],
			},
			"mdm_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["mdm_url"],
			},
			"service_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"org_admin_username"},
				RequiredWith:  []string{"service_private_key"},
				DefaultFunc:   schema.EnvDefaultFunc(ServiceID, nil),
				Description:   descriptions["service_id"],
			},
			"service_private_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"org_admin_password"},
				RequiredWith:  []string{"service_id"},
				DefaultFunc:   schema.EnvDefaultFunc(ServicePK, nil),
				Description:   descriptions["service_private_key"],
			},
			"oauth2_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ClientID, nil),
				Description: descriptions["oauth2_client_id"],
			},
			"oauth2_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(ClientPassword, nil),
				Description: descriptions["oauth2_password"],
			},
			"org_admin_username": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   descriptions["org_admin_username"],
				RequiredWith:  []string{"org_admin_password"},
				ConflictsWith: []string{"service_id"},
				DefaultFunc:   schema.EnvDefaultFunc(OrgAdminUsername, nil),
			},
			"org_admin_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   descriptions["org_admin_password"],
				RequiredWith:  []string{"org_admin_username"},
				ConflictsWith: []string{"service_private_key"},
				DefaultFunc:   schema.EnvDefaultFunc(OrgAdminPassword, nil),
			},
			"uaa_username": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["uaa_username"],
				RequiredWith: []string{"uaa_password"},
				DefaultFunc:  schema.EnvDefaultFunc(UAAUsername, nil),
			},
			"uaa_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Description:  descriptions["uaa_password"],
				RequiredWith: []string{"uaa_username"},
				DefaultFunc:  schema.EnvDefaultFunc(UAAPassword, nil),
			},
			"uaa_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["uaa_url"],
			},
			"shared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   false,
				DefaultFunc: schema.EnvDefaultFunc(SharedKey, nil),
				Description: descriptions["shared_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(SecretKey, nil),
				Description: descriptions["secret_key"],
			},
			"retry_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: descriptions["retry_max"],
			},
			"debug_log": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(DebugLog, nil),
				Description: descriptions["debug_log"],
			},
			"debug_stderr": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(DebugStdErr, nil),
				Description: descriptions["debug_stderr"],
			},
			"credentials": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hsdp_iam_org":                                   organization.ResourceIAMOrg(),
			"hsdp_iam_group":                                 group.ResourceIAMGroup(),
			"hsdp_iam_role":                                  role.ResourceIAMRole(),
			"hsdp_iam_proposition":                           proposition.ResourceIAMProposition(),
			"hsdp_iam_application":                           application.ResourceIAMApplication(),
			"hsdp_iam_user":                                  user.ResourceIAMUser(),
			"hsdp_iam_client":                                client.ResourceIAMClient(),
			"hsdp_iam_service":                               service.ResourceIAMService(),
			"hsdp_iam_mfa_policy":                            iam.ResourceIAMMFAPolicy(),
			"hsdp_iam_password_policy":                       iam.ResourceIAMPasswordPolicy(),
			"hsdp_iam_email_template":                        email_template.ResourceIAMEmailTemplate(),
			"hsdp_edge_app":                                  edge.ResourceEdgeApp(),
			"hsdp_edge_config":                               edge.ResourceEdgeConfig(),
			"hsdp_edge_custom_cert":                          edge.ResourceEdgeCustomCert(),
			"hsdp_edge_sync":                                 edge.ResourceEdgeSync(),
			"hsdp_iam_sms_gateway":                           iam.ResourceIAMSMSGatewayConfig(),
			"hsdp_iam_sms_template":                          iam.ResourceIAMSMSTemplate(),
			"hsdp_iam_activation_email":                      iam.ResourceIAMActivationEmail(),
			"hsdp_connect_mdm_standard_service":              mdm.ResourceConnectMDMStandardService(),
			"hsdp_connect_mdm_service_action":                mdm.ResourceConnectMDMServiceAction(),
			"hsdp_connect_mdm_device_group":                  mdm.ResourceConnectMDMDeviceGroup(),
			"hsdp_connect_mdm_device_type":                   mdm.ResourceConnectMDMDeviceType(),
			"hsdp_connect_mdm_oauth_client":                  mdm.ResourceConnectMDMOAuthClient(),
			"hsdp_connect_mdm_authentication_method":         mdm.ResourceConnectMDMAuthenticationMethod(),
			"hsdp_connect_mdm_service_reference":             mdm.ResourceConnectMDMServiceReference(),
			"hsdp_connect_mdm_bucket":                        mdm.ResourceConnectMDMBucket(),
			"hsdp_connect_mdm_data_type":                     mdm.ResourceConnectMDMDataType(),
			"hsdp_connect_mdm_blob_data_contract":            mdm.ResourceConnectMDMBlobDataContract(),
			"hsdp_connect_mdm_blob_subscription":             mdm.ResourceConnectMDMBlobSubscription(),
			"hsdp_connect_mdm_firmware_component":            mdm.ResourceConnectMDMFirmwareComponent(),
			"hsdp_connect_mdm_proposition":                   mdm.ResourceMDMProposition(),
			"hsdp_connect_mdm_application":                   mdm.ResourceMDMApplication(),
			"hsdp_connect_mdm_firmware_component_version":    mdm.ResourceConnectMDMFirmwareComponentVersion(),
			"hsdp_connect_mdm_firmware_distribution_request": mdm.ResourceConnectMDMFirmwareDistributionRequest(),
			"hsdp_connect_iot_provisioning_orgconfiguration": provisioning.ResourceConnectIoTProvisioningOrgConfiguration(),
			"hsdp_iam_group_membership":                      group_membership.ResourceIAMGroupMembership(),
			"hsdp_iam_role_sharing_policy":                   role_sharing_policy.ResourceRoleSharingPolicy(),
			"hsdp_iam_device":                                device.ResourceIAMDevice(),
			"hsdp_blr_bucket":                                blr.ResourceBLRBucket(),
			"hsdp_blr_blob_store_policy":                     blr.ResourceBLRBlobStorePolicy(),
			"hsdp_dbs_sqs_subscriber":                        dbs.ResourceDBSSQSSubscriber(),
			"hsdp_dbs_topic_subscription":                    dbs.ResourceDBSTopicSubscription(),
			"hsdp_tenant_key":                                tenant.ResourceTenantKey(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hsdp_iam_introspect":                            iam.DataSourceIAMIntrospect(),
			"hsdp_iam_user":                                  user.DataSourceUser(),
			"hsdp_iam_service":                               service.DataSourceService(),
			"hsdp_iam_permissions":                           iam.DataSourceIAMPermissions(),
			"hsdp_iam_org":                                   organization.DataSourceIAMOrg(),
			"hsdp_iam_proposition":                           proposition.DataSourceIAMProposition(),
			"hsdp_iam_application":                           application.DataSourceIAMApplication(),
			"hsdp_config":                                    configuration.DataSourceConfig(),
			"hsdp_edge_device":                               edge.DataSourceEdgeDevice(),
			"hsdp_iam_group":                                 group.DataSourceIAMGroup(),
			"hsdp_iam_role":                                  role.DataSourceIAMRole(),
			"hsdp_iam_users":                                 user.DataSourceIAMUsers(),
			"hsdp_iam_client":                                client.DataSourceIAMClient(),
			"hsdp_connect_mdm_proposition":                   mdm.DataSourceConnectMDMProposition(),
			"hsdp_connect_mdm_application":                   mdm.DataSourceConnectMDMApplication(),
			"hsdp_connect_mdm_standard_services":             mdm.DataSourceConnectMDMStandardServices(),
			"hsdp_connect_mdm_regions":                       mdm.DataSourceConnectMDMRegions(),
			"hsdp_connect_mdm_oauth_client_scopes":           mdm.DataSourceConnectMDMOauthClientScopes(),
			"hsdp_connect_mdm_region":                        mdm.DataSourceConnectMDMRegion(),
			"hsdp_connect_mdm_resource_limits":               mdm.DataSourceResourceLimits(),
			"hsdp_connect_mdm_subscriber_types":              mdm.DataSourceConnectMDMSubscriberTypes(),
			"hsdp_connect_mdm_storage_classes":               mdm.DataSourceConnectMDMStorageClasses(),
			"hsdp_connect_mdm_storage_class":                 mdm.DataSourceConnectMDMStorageClass(),
			"hsdp_connect_mdm_standard_service":              mdm.DataSourceConnectMDMStandardService(),
			"hsdp_connect_mdm_data_subscribers":              mdm.DataSourceConnectMDMDataSubscribers(),
			"hsdp_connect_mdm_data_adapters":                 mdm.DataSourceConnectMDMDataAdapters(),
			"hsdp_connect_iot_provisioning_orgconfiguration": provisioning.DataSourceConnectIoTProvisioningOrgConfiguration(),
			"hsdp_iam_email_templates":                       email_template.DataSourceIAMEmailTemplates(),
			"hsdp_connect_mdm_bucket":                        mdm.DataSourceConnectMDMBucket(),
			"hsdp_connect_mdm_data_type":                     mdm.DataSourceConnectMDMDataType(),
			"hsdp_iam_token":                                 iam.DataSourceIAMToken(),
			"hsdp_connect_mdm_service_agent":                 mdm.DataSourceConnectMDMServiceAgent(),
			"hsdp_connect_mdm_service_agents":                mdm.DataSourceConnectMDMServiceAgents(),
			"hsdp_iam_permission":                            iam.DataSourceIAMPermission(),
			"hsdp_iam_role_sharing_policies":                 role_sharing_policy.DataSourceIAMRoleSharingPolicies(),
			"hsdp_discovery_service":                         discovery.DataSourceDiscoveryService(),
			"hsdp_connect_mdm_service_action":                mdm.DataSourceConnectMDMServiceAction(),
			"hsdp_connect_mdm_service_actions":               mdm.DataSourceConnectMDMServiceActions(),
			"hsdp_blr_store_policy":                          blr.DataSourceBLRBlobStorePolicyDefinition(),
		},
		ConfigureContextFunc: providerConfigure(build),
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"region":              "The HSDP region to configure for",
		"environment":         "The HSDP environment to configure for",
		"iam_url":             "The HSDP IAM instance URL",
		"idm_url":             "The HSDP IDM instance URL",
		"mdm_url":             "The Connect MDM URL to use",
		"oauth2_client_id":    "The OAuth2 client id",
		"oauth2_password":     "The OAuth2 password",
		"service_id":          "The service ID to use as Organization Admin",
		"service_private_key": "The private key of the service ID",
		"org_admin_username":  "The username of the Organization Admin",
		"org_admin_password":  "The password of the Organization Admin",
		"shared_key":          "The shared key",
		"secret_key":          "The secret key",
		"debug_log":           "The log file to write debugging output to",
		"debug_stderr":        "Debug to stderr",

		"retry_max":           "Maximum number of retries for API requests",
		"uaa_username":        "The username of the Cloudfoundry account to use",
		"uaa_password":        "The password of the Cloudfoundry account to use",
		"uaa_url":             "The URL of the UAA server",
	}
}

func providerConfigure(build string) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics

		c := &config.Config{}

		c.BuildVersion = build
		c.Region = d.Get("region").(string)
		c.Environment = d.Get("environment").(string)
		c.IAMURL = d.Get("iam_url").(string)
		c.IDMURL = d.Get("idm_url").(string)
		c.OAuth2ClientID = d.Get("oauth2_client_id").(string)
		c.OAuth2ClientSecret = d.Get("oauth2_password").(string)
		c.ServiceID = d.Get("service_id").(string)
		c.ServicePrivateKey = d.Get("service_private_key").(string)
		c.OrgAdminUsername = d.Get("org_admin_username").(string)
		c.OrgAdminPassword = d.Get("org_admin_password").(string)
		c.SharedKey = d.Get("shared_key").(string)
		c.SecretKey = d.Get("secret_key").(string)
		c.DebugLog = d.Get("debug_log").(string)
		c.DebugStdErr = d.Get("debug_stderr").(bool)
		c.RetryMax = d.Get("retry_max").(int)
		c.UAAUsername = d.Get("uaa_username").(string)
		c.UAAPassword = d.Get("uaa_password").(string)
		c.UAAURL = d.Get("uaa_url").(string)
		c.TimeZone = "UTC"
		c.MDMURL = d.Get("mdm_url").(string)

		credentialsFile := d.Get("credentials").(string)
		if credentialsFile != "" {
			file, _ := os.ReadFile(credentialsFile)
			_ = json.Unmarshal(file, &c)
		}
		if c.DebugLog != "" {
			debugFile, err := os.OpenFile(c.DebugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if err == nil {
				c.DebugWriter = debugFile
			}
		}
		if c.DebugStdErr && c.DebugWriter == nil { // Crossplane
			c.DebugWriter = os.Stderr
		}
		c.SetupIAMClient()
		c.SetupSTLClient()
		c.SetupMDMClient()
		c.SetupDiscoveryClient()
		c.SetupBLRClient()
		c.SetupDBSClient()
		c.SetupProvisioningClient()

		ma, err := jsonformat.NewMarshaller(false, "", "", fhirversion.STU3)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		c.STU3MA = ma

		um, err := jsonformat.NewUnmarshaller("UTC", fhirversion.STU3)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		c.STU3UM = um

		ma, err = jsonformat.NewMarshaller(false, "", "", fhirversion.R4)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		c.R4MA = ma

		um, err = jsonformat.NewUnmarshaller("UTC", fhirversion.R4)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		c.R4UM = um

		return c, diags
	}
}
