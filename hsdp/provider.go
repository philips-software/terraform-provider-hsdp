package hsdp

import (
	"context"
	"encoding/json"
	"os"

	"github.com/philips-software/terraform-provider-hsdp/internal/services/connect/dbs"

	"github.com/philips-software/terraform-provider-hsdp/internal/services/blr"

	"github.com/philips-software/terraform-provider-hsdp/internal/services/iam/group_membership"

	"github.com/google/fhir/go/fhirversion"
	"github.com/google/fhir/go/jsonformat"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/ai/inference"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/ai/workspace"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/cdl"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/cdr/fhir_store"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/cdr/org"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/cdr/practitioner"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/cdr/subscription"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/ch"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/configuration"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/dicom"
	repository2 "github.com/philips-software/terraform-provider-hsdp/internal/services/dicom/repository"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/discovery"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/docker/namespace"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/docker/repository"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/docker/service_key"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/edge"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/function"
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
	"github.com/philips-software/terraform-provider-hsdp/internal/services/metrics"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/notification"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/pki"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/pki/tenant"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/s3creds"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

const (
	Region           = "HSDP_REGION"
	Environment      = "HSDP_ENVIRONMENT"
	CartelSecret     = "HSDP_CARTEL_SECRET"
	CartelToken      = "HSDP_CARTEL_TOKEN"
	CartelHost       = "HSDP_CARTEL_HOST"
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
			"cartel_host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(CartelHost, nil),
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
				DefaultFunc: schema.EnvDefaultFunc(DebugLog, nil),
				Description: descriptions["debug_log"],
			},
			"debug_stderr": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(DebugStdErr, nil),
				Description: descriptions["debug_stderr"],
			},
			"ai_inference_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
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
			"hsdp_s3creds_policy":                            s3creds.ResourceS3CredsPolicy(),
			"hsdp_container_host":                            ch.ResourceContainerHost(),
			"hsdp_metrics_autoscaler":                        metrics.ResourceMetricsAutoscaler(),
			"hsdp_cdr_org":                                   org.ResourceCDROrg(),
			"hsdp_cdr_subscription":                          subscription.ResourceCDRSubscription(),
			"hsdp_dicom_store_config":                        dicom.ResourceDICOMStoreConfig(),
			"hsdp_dicom_object_store":                        dicom.ResourceDICOMObjectStore(),
			"hsdp_dicom_repository":                          repository2.ResourceDICOMRepository(),
			"hsdp_pki_tenant":                                tenant.ResourcePKITenant(),
			"hsdp_pki_cert":                                  pki.ResourcePKICert(),
			"hsdp_edge_app":                                  edge.ResourceEdgeApp(),
			"hsdp_edge_config":                               edge.ResourceEdgeConfig(),
			"hsdp_edge_custom_cert":                          edge.ResourceEdgeCustomCert(),
			"hsdp_edge_sync":                                 edge.ResourceEdgeSync(),
			"hsdp_function":                                  function.ResourceFunction(),
			"hsdp_notification_producer":                     notification.ResourceNotificationProducer(),
			"hsdp_notification_subscriber":                   notification.ResourceNotificationSubscriber(),
			"hsdp_notification_topic":                        notification.ResourceNotificationTopic(),
			"hsdp_notification_subscription":                 notification.ResourceNotificationSubscription(),
			"hsdp_ai_inference_compute_environment":          inference.ResourceAIInferenceComputeEnvironment(),
			"hsdp_ai_inference_compute_target":               inference.ResourceAIInferenceComputeTarget(),
			"hsdp_ai_inference_model":                        inference.ResourceAIInferenceModel(),
			"hsdp_ai_inference_job":                          inference.ResourceAIInferenceJob(),
			"hsdp_dicom_gateway_config":                      dicom.ResourceDICOMGatewayConfig(),
			"hsdp_cdl_research_study":                        cdl.ResourceCDLResearchStudy(),
			"hsdp_dicom_remote_node":                         dicom.ResourceDICOMRemoteNode(),
			"hsdp_cdl_data_type_definition":                  cdl.ResourceCDLDataTypeDefinition(),
			"hsdp_cdl_label_definition":                      cdl.ResourceCDLLabelDefinition(),
			"hsdp_cdl_export_route":                          cdl.ResourceCDLExportRoute(),
			"hsdp_ai_workspace_compute_target":               workspace.ResourceAIWorkspaceComputeTarget(),
			"hsdp_ai_workspace":                              workspace.ResourceAIWorkspace(),
			"hsdp_iam_sms_gateway":                           iam.ResourceIAMSMSGatewayConfig(),
			"hsdp_iam_sms_template":                          iam.ResourceIAMSMSTemplate(),
			"hsdp_iam_activation_email":                      iam.ResourceIAMActivationEmail(),
			"hsdp_docker_service_key":                        service_key.ResourceDockerServiceKey(),
			"hsdp_docker_namespace":                          namespace.ResourceDockerNamespace(),
			"hsdp_docker_namespace_user":                     namespace.ResourceDockerNamespaceUser(),
			"hsdp_docker_repository":                         repository.ResourceDockerRepository(),
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
			"hsdp_iam_group_membership":                      group_membership.ResourceIAMGroupMembership(),
			"hsdp_dicom_notification":                        dicom.ResourceDICOMNotification(),
			"hsdp_cdr_practitioner":                          practitioner.ResourceCDRPractitioner(),
			"hsdp_iam_role_sharing_policy":                   role_sharing_policy.ResourceRoleSharingPolicy(),
			"hsdp_iam_device":                                device.ResourceIAMDevice(),
			"hsdp_blr_bucket":                                blr.ResourceBLRBucket(),
			"hsdp_blr_blob_store_policy":                     blr.ResourceBLRBlobStorePolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hsdp_iam_introspect":                        iam.DataSourceIAMIntrospect(),
			"hsdp_iam_user":                              user.DataSourceUser(),
			"hsdp_iam_service":                           service.DataSourceService(),
			"hsdp_iam_permissions":                       iam.DataSourceIAMPermissions(),
			"hsdp_iam_org":                               organization.DataSourceIAMOrg(),
			"hsdp_iam_proposition":                       proposition.DataSourceIAMProposition(),
			"hsdp_iam_application":                       application.DataSourceIAMApplication(),
			"hsdp_s3creds_access":                        s3creds.DataSourceS3CredsAccess(),
			"hsdp_s3creds_policy":                        s3creds.DataSourceS3CredsPolicy(),
			"hsdp_config":                                configuration.DataSourceConfig(),
			"hsdp_container_host_subnet_types":           ch.DataSourceContainerHostSubnetTypes(),
			"hsdp_cdr_fhir_store":                        fhir_store.DataSourceCDRFHIRStore(),
			"hsdp_pki_root":                              pki.DataSourcePKIRoot(),
			"hsdp_pki_policy":                            pki.DataSourcePKIPolicy(),
			"hsdp_edge_device":                           edge.DataSourceEdgeDevice(),
			"hsdp_notification_producers":                notification.DataSourceNotificationProducers(),
			"hsdp_notification_producer":                 notification.DataSourceNotificationProducer(),
			"hsdp_notification_topics":                   notification.DataSourceNotificationTopics(),
			"hsdp_notification_topic":                    notification.DataSourceNotificationTopic(),
			"hsdp_notification_subscription":             notification.DataSourceNotificationSubscription(),
			"hsdp_notification_subscriber":               notification.DataSourceNotificationSubscriber(),
			"hsdp_ai_inference_service_instance":         inference.DataSourceAIInferenceServiceInstance(),
			"hsdp_ai_inference_compute_environments":     inference.DataSourceAIInferenceComputeEnvironments(),
			"hsdp_ai_inference_compute_targets":          inference.DataSourceAIInferenceComputeTargets(),
			"hsdp_ai_inference_jobs":                     inference.DataSourceAIInferenceJobs(),
			"hsdp_ai_inference_models":                   inference.DataSourceAIInferenceModels(),
			"hsdp_cdl_instance":                          cdl.DataSourceCDLInstance(),
			"hsdp_cdl_research_study":                    cdl.DataSourceCDLResearchStudy(),
			"hsdp_cdl_research_studies":                  cdl.DataSourceCDLResearchStudies(),
			"hsdp_container_host_instances":              ch.DataSourceContainerHostInstances(),
			"hsdp_cdl_data_type_definitions":             cdl.DataSourceCDLDataTypeDefinitions(),
			"hsdp_cdl_data_type_definition":              cdl.DataSourceCDLDataTypeDefinition(),
			"hsdp_cdl_label_definition":                  cdl.DataSourceCDLLabelDefinition(),
			"hsdp_cdl_export_route":                      cdl.DataSourceCDLExportRoute(),
			"hsdp_ai_workspace_service_instance":         workspace.DataSourceAIWorkspaceServiceInstance(),
			"hsdp_ai_workspace_compute_targets":          workspace.DataSourceAIWorkspaceComputeTargets(),
			"hsdp_ai_workspace":                          workspace.DataSourceAIWorkspace(),
			"hsdp_iam_group":                             group.DataSourceIAMGroup(),
			"hsdp_iam_role":                              role.DataSourceIAMRole(),
			"hsdp_iam_users":                             user.DataSourceIAMUsers(),
			"hsdp_docker_namespace":                      namespace.DataSourceDockerNamespace(),
			"hsdp_docker_namespaces":                     namespace.DataSourceDockerNamespaces(),
			"hsdp_docker_repository":                     repository.DataSourceDockerRepository(),
			"hsdp_iam_client":                            client.DataSourceIAMClient(),
			"hsdp_connect_mdm_proposition":               mdm.DataSourceConnectMDMProposition(),
			"hsdp_connect_mdm_application":               mdm.DataSourceConnectMDMApplication(),
			"hsdp_connect_mdm_standard_services":         mdm.DataSourceConnectMDMStandardServices(),
			"hsdp_connect_mdm_regions":                   mdm.DataSourceConnectMDMRegions(),
			"hsdp_connect_mdm_oauth_client_scopes":       mdm.DataSourceConnectMDMOauthClientScopes(),
			"hsdp_connect_mdm_region":                    mdm.DataSourceConnectMDMRegion(),
			"hsdp_connect_mdm_resource_limits":           mdm.DataSourceResourceLimits(),
			"hsdp_connect_mdm_subscriber_types":          mdm.DataSourceConnectMDMSubscriberTypes(),
			"hsdp_connect_mdm_storage_classes":           mdm.DataSourceConnectMDMStorageClasses(),
			"hsdp_connect_mdm_storage_class":             mdm.DataSourceConnectMDMStorageClass(),
			"hsdp_connect_mdm_standard_service":          mdm.DataSourceConnectMDMStandardService(),
			"hsdp_connect_mdm_data_subscribers":          mdm.DataSourceConnectMDMDataSubscribers(),
			"hsdp_connect_mdm_data_adapters":             mdm.DataSourceConnectMDMDataAdapters(),
			"hsdp_iam_email_templates":                   email_template.DataSourceIAMEmailTemplates(),
			"hsdp_connect_mdm_bucket":                    mdm.DataSourceConnectMDMBucket(),
			"hsdp_connect_mdm_data_type":                 mdm.DataSourceConnectMDMDataType(),
			"hsdp_container_host_security_groups":        ch.DataSourceContainerHostSecurityGroups(),
			"hsdp_container_host_security_group_details": ch.DataSourceContainerHostSecurityGroupDetails(),
			"hsdp_iam_token":                             iam.DataSourceIAMToken(),
			"hsdp_connect_mdm_service_agent":             mdm.DataSourceConnectMDMServiceAgent(),
			"hsdp_connect_mdm_service_agents":            mdm.DataSourceConnectMDMServiceAgents(),
			"hsdp_container_host":                        ch.DataSourceContainerHost(),
			"hsdp_iam_permission":                        iam.DataSourceIAMPermission(),
			"hsdp_cdr_practitioner":                      practitioner.DataSourceCDRPractitioner(),
			"hsdp_cdr_org":                               org.DataSourceCDROrg(),
			"hsdp_iam_role_sharing_policies":             role_sharing_policy.DataSourceIAMRoleSharingPolicies(),
			"hsdp_discovery_service":                     discovery.DataSourceDiscoveryService(),
			"hsdp_connect_mdm_service_action":            mdm.DataSourceConnectMDMServiceAction(),
			"hsdp_connect_mdm_service_actions":           mdm.DataSourceConnectMDMServiceActions(),
			"hsdp_blr_store_policy":                      blr.DataSourceBLRBlobStorePolicyDefinition(),
			"hsdp_dbs_topic_subscription":                dbs.ResourceDBSTopicSubscription(),
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
		c.S3CredsURL = d.Get("s3creds_url").(string)
		c.CartelHost = d.Get("cartel_host").(string)
		c.CartelToken = d.Get("cartel_token").(string)
		c.CartelSecret = d.Get("cartel_secret").(string)
		c.CartelNoTLS = d.Get("cartel_no_tls").(bool)
		c.DebugStdErr = d.Get("debug_stderr").(bool)
		c.CartelSkipVerify = d.Get("cartel_skip_verify").(bool)
		c.RetryMax = d.Get("retry_max").(int)
		c.UAAUsername = d.Get("uaa_username").(string)
		c.UAAPassword = d.Get("uaa_password").(string)
		c.UAAURL = d.Get("uaa_url").(string)
		c.NotificationURL = d.Get("notification_url").(string)
		c.TimeZone = "UTC"
		c.AIInferenceEndpoint = d.Get("ai_inference_endpoint").(string)
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
		c.SetupS3CredsClient()
		c.SetupCartelClient()
		c.SetupConsoleClient()
		c.SetupPKIClient()
		c.SetupSTLClient()
		c.SetupNotificationClient()
		c.SetupMDMClient()
		c.SetupDiscoveryClient()
		c.SetupBLRClient()
		c.SetupDBSClient()

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
