package provider

import (
	"context"
	"os"

	"github.com/google/fhir/go/jsonformat"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	config2 "github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/ai"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/cdl"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/cdr"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/ch"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/dicom"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/edge"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/function"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/metrics"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/notification"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/pki"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/s3creds"
)

const (
	Region           = "HSDP_REGION"
	Environment      = "HSDP_ENVIRONMENT"
	CartelSecret     = "HSDP_CARTEL_SECRET"
	CartelToken      = "HSDP_CARTEL_TOKEN"
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
)

// Provider returns an instance of the HSDP provider
func Provider(build string) *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(Region, nil),
				Description: descriptions["region"],
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(Environment, "client-test"),
				Description: descriptions["environment"],
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
			"s3creds_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["s3creds_url"],
			},
			"notification_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["notification_url"],
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
			"cartel_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["cartel_host"],
			},
			"cartel_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(CartelToken, nil),
				Description: descriptions["cartel_token"],
			},
			"cartel_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(CartelSecret, nil),
				Description: descriptions["cartel_secret"],
			},
			"cartel_no_tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["cartel_no_tls"],
			},
			"cartel_skip_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: descriptions["cartel_skip_verify"],
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
				Description: descriptions["debug_log"],
			},
			"ai_inference_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hsdp_iam_org":                          iam.ResourceIAMOrg(),
			"hsdp_iam_group":                        iam.ResourceIAMGroup(),
			"hsdp_iam_role":                         iam.ResourceIAMRole(),
			"hsdp_iam_proposition":                  iam.ResourceIAMProposition(),
			"hsdp_iam_application":                  iam.ResourceIAMApplication(),
			"hsdp_iam_user":                         iam.ResourceIAMUser(),
			"hsdp_iam_client":                       iam.ResourceIAMClient(),
			"hsdp_iam_service":                      iam.ResourceIAMService(),
			"hsdp_iam_mfa_policy":                   iam.ResourceIAMMFAPolicy(),
			"hsdp_iam_password_policy":              iam.ResourceIAMPasswordPolicy(),
			"hsdp_iam_email_template":               iam.ResourceIAMEmailTemplate(),
			"hsdp_s3creds_policy":                   s3creds.ResourceS3CredsPolicy(),
			"hsdp_container_host":                   ch.ResourceContainerHost(),
			"hsdp_container_host_exec":              ch.ResourceContainerHostExec(),
			"hsdp_metrics_autoscaler":               metrics.ResourceMetricsAutoscaler(),
			"hsdp_cdr_org":                          cdr.ResourceCDROrg(),
			"hsdp_cdr_subscription":                 cdr.ResourceCDRSubscription(),
			"hsdp_dicom_store_config":               dicom.ResourceDICOMStoreConfig(),
			"hsdp_dicom_object_store":               dicom.ResourceDICOMObjectStore(),
			"hsdp_dicom_repository":                 dicom.ResourceDICOMRepository(),
			"hsdp_pki_tenant":                       pki.ResourcePKITenant(),
			"hsdp_pki_cert":                         pki.ResourcePKICert(),
			"hsdp_edge_app":                         edge.ResourceEdgeApp(),
			"hsdp_edge_config":                      edge.ResourceEdgeConfig(),
			"hsdp_edge_custom_cert":                 edge.ResourceEdgeCustomCert(),
			"hsdp_edge_sync":                        edge.ResourceEdgeSync(),
			"hsdp_function":                         function.ResourceFunction(),
			"hsdp_notification_producer":            notification.ResourceNotificationProducer(),
			"hsdp_notification_subscriber":          notification.ResourceNotificationSubscriber(),
			"hsdp_notification_topic":               notification.ResourceNotificationTopic(),
			"hsdp_notification_subscription":        notification.ResourceNotificationSubscription(),
			"hsdp_ai_inference_compute_environment": ai.ResourceAIInferenceComputeEnvironment(),
			"hsdp_ai_inference_compute_target":      ai.ResourceAIInferenceComputeTarget(),
			"hsdp_ai_inference_model":               ai.ResourceAIInferenceModel(),
			"hsdp_ai_inference_job":                 ai.ResourceAIInferenceJob(),
			"hsdp_dicom_gateway_config":             dicom.ResourceDICOMGatewayConfig(),
			"hsdp_cdl_research_study":               cdl.ResourceCDLResearchStudy(),
			"hsdp_dicom_remote_node":                dicom.ResourceDICOMRemoteNode(),
			"hsdp_cdl_data_type_definition":         cdl.ResourceCDLDataTypeDefinition(),
			"hsdp_cdl_label_definition":             cdl.ResourceCDLLabelDefinition(),
			"hsdp_cdl_export_route":                 cdl.ResourceCDLExportRoute(),
			"hsdp_ai_workspace_compute_target":      ai.ResourceAIWorkspaceComputeTarget(),
			"hsdp_ai_workspace":                     ai.ResourceAIWorkspace(),
			"hsdp_iam_sms_gateway":                  iam.ResourceIAMSMSGatewayConfig(),
			"hsdp_iam_sms_template":                 iam.ResourceIAMSMSTemplate(),
			"hsdp_iam_activation_email":             iam.ResourceIAMActivationEmail(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hsdp_iam_introspect":                    iam.DataSourceIAMIntrospect(),
			"hsdp_iam_user":                          iam.DataSourceUser(),
			"hsdp_iam_service":                       iam.DataSourceService(),
			"hsdp_iam_permissions":                   iam.DataSourceIAMPermissions(),
			"hsdp_iam_org":                           iam.DataSourceIAMOrg(),
			"hsdp_iam_proposition":                   iam.DataSourceIAMProposition(),
			"hsdp_iam_application":                   iam.DataSourceIAMApplication(),
			"hsdp_s3creds_access":                    s3creds.DataSourceS3CredsAccess(),
			"hsdp_s3creds_policy":                    s3creds.DataSourceS3CredsPolicy(),
			"hsdp_config":                            config.DataSourceConfig(),
			"hsdp_container_host_subnet_types":       ch.DataSourceContainerHostSubnetTypes(),
			"hsdp_cdr_fhir_store":                    cdr.DataSourceCDRFHIRStore(),
			"hsdp_pki_root":                          pki.DataSourcePKIRoot(),
			"hsdp_pki_policy":                        pki.DataSourcePKIPolicy(),
			"hsdp_edge_device":                       edge.DataSourceEdgeDevice(),
			"hsdp_notification_producers":            notification.DataSourceNotificationProducers(),
			"hsdp_notification_producer":             notification.DataSourceNotificationProducer(),
			"hsdp_notification_topics":               notification.DataSourceNotificationTopics(),
			"hsdp_notification_topic":                notification.DataSourceNotificationTopic(),
			"hsdp_notification_subscription":         notification.DataSourceNotificationSubscription(),
			"hsdp_notification_subscriber":           notification.DataSourceNotificationSubscriber(),
			"hsdp_ai_inference_service_instance":     ai.DataSourceAIInferenceServiceInstance(),
			"hsdp_ai_inference_compute_environments": ai.DataSourceAIInferenceComputeEnvironments(),
			"hsdp_ai_inference_compute_targets":      ai.DataSourceAIInferenceComputeTargets(),
			"hsdp_ai_inference_jobs":                 ai.DataSourceAIInferenceJobs(),
			"hsdp_ai_inference_models":               ai.DataSourceAIInferenceModels(),
			"hsdp_cdl_instance":                      cdl.DataSourceCDLInstance(),
			"hsdp_cdl_research_study":                cdl.DataSourceCDLResearchStudy(),
			"hsdp_cdl_research_studies":              cdl.DataSourceCDLResearchStudies(),
			"hsdp_container_host_instances":          ch.DataSourceContainerHostInstances(),
			"hsdp_cdl_data_type_definitions":         cdl.DataSourceCDLDataTypeDefinitions(),
			"hsdp_cdl_data_type_definition":          cdl.DataSourceCDLDataTypeDefinition(),
			"hsdp_cdl_label_definition":              cdl.DataSourceCDLLabelDefinition(),
			"hsdp_cdl_export_route":                  cdl.DataSourceCDLExportRoute(),
			"hsdp_ai_workspace_service_instance":     ai.DataSourceAIWorkspaceServiceInstance(),
			"hsdp_ai_workspace_compute_targets":      ai.DataSourceAIWorkspaceComputeTargets(),
			"hsdp_ai_workspace":                      ai.DataSourceAIWorkspace(),
			"hsdp_iam_group":                         iam.DataSourceIAMGroup(),
			"hsdp_iam_role":                          iam.DataSourceIAMRole(),
			"hsdp_iam_users":                         iam.DataSourceIAMUsers(),
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
		"s3creds_url":         "The HSDP S3 Credentials instance URL",
		"notification_url":    "The HSDP Notification service base URL to use",
		"oauth2_client_id":    "The OAuth2 client id",
		"oauth2_password":     "The OAuth2 password",
		"service_id":          "The service ID to use as Organization Admin",
		"service_private_key": "The private key of the service ID",
		"org_admin_username":  "The username of the Organization Admin",
		"org_admin_password":  "The password of the Organization Admin",
		"shared_key":          "The shared key",
		"secret_key":          "The secret key",
		"debug_log":           "The log file to write debugging output to",
		"cartel_host":         "The Cartel host",
		"cartel_token":        "The Cartel token key",
		"cartel_secret":       "The Cartel secret key",
		"cartel_no_tls":       "Disable TLS for Cartel",
		"cartel_skip_verify":  "Skip certificate verification",
		"retry_max":           "Maximum number of retries for API requests",
		"uaa_username":        "The username of the Cloudfoundry account to use",
		"uaa_password":        "The password of the Cloudfoundry account to use",
		"uaa_url":             "The URL of the UAA server",
	}
}

func providerConfigure(build string) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics

		config := &config2.Config{}

		config.BuildVersion = build
		config.Region = d.Get("region").(string)
		config.Environment = d.Get("environment").(string)
		config.IAMURL = d.Get("iam_url").(string)
		config.IDMURL = d.Get("idm_url").(string)
		config.OAuth2ClientID = d.Get("oauth2_client_id").(string)
		config.OAuth2Secret = d.Get("oauth2_password").(string)
		config.ServiceID = d.Get("service_id").(string)
		config.ServicePrivateKey = d.Get("service_private_key").(string)
		config.OrgAdminUsername = d.Get("org_admin_username").(string)
		config.OrgAdminPassword = d.Get("org_admin_password").(string)
		config.SharedKey = d.Get("shared_key").(string)
		config.SecretKey = d.Get("secret_key").(string)
		config.DebugLog = d.Get("debug_log").(string)
		config.S3CredsURL = d.Get("s3creds_url").(string)
		config.CartelHost = d.Get("cartel_host").(string)
		config.CartelToken = d.Get("cartel_token").(string)
		config.CartelSecret = d.Get("cartel_secret").(string)
		config.CartelNoTLS = d.Get("cartel_no_tls").(bool)
		config.CartelSkipVerify = d.Get("cartel_skip_verify").(bool)
		config.RetryMax = d.Get("retry_max").(int)
		config.UAAUsername = d.Get("uaa_username").(string)
		config.UAAPassword = d.Get("uaa_password").(string)
		config.UAAURL = d.Get("uaa_url").(string)
		config.NotificationURL = d.Get("notification_url").(string)
		config.TimeZone = "UTC"
		config.AIInferenceEndpoint = d.Get("ai_inference_endpoint").(string)

		config.SetupIAMClient()
		config.SetupS3CredsClient()
		config.SetupCartelClient()
		config.SetupConsoleClient()
		config.SetupPKIClient()
		config.SetupSTLClient()
		config.SetupNotificationClient()

		if config.DebugLog != "" {
			debugFile, err := os.OpenFile(config.DebugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				config.DebugFile = nil
			} else {
				config.DebugFile = debugFile
			}
		}

		ma, err := jsonformat.NewMarshaller(false, "", "", jsonformat.STU3)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		config.Ma = ma

		um, err := jsonformat.NewUnmarshaller("UTC", jsonformat.STU3)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		config.Um = um

		return config, diags
	}
}
