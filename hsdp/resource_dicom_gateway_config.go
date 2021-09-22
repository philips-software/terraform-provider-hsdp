package hsdp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/dicom"
)

func schemaApplicationEntity() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		MaxItems: 100,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allow_any": {
					Type:     schema.TypeBool,
					Required: true,
				},
				"ae_title": {
					Type:     schema.TypeString,
					Required: true,
				},
				"organization_id": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func resourceDICOMGatewayConfig() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceDICOMGatewayConfigCreate,
		ReadContext:   resourceDICOMGatewayConfigRead,
		DeleteContext: resourceDICOMGatewayConfigDelete,

		Schema: map[string]*schema.Schema{
			"config_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"store_service": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"is_secure": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						// ---Advanced features start
						"pdu_length": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  65535,
						},
						"artim_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3000,
						},
						"association_idle_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  4500,
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"authenticate_client_certificate": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						// ---Advanced features end
						"application_entity": schemaApplicationEntity(),
					},
				},
			},
			"query_retrieve_service": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"is_secure": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						// ---Advanced features start
						"pdu_length": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  65535,
						},
						"artim_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3000,
						},
						"association_idle_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  4500,
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"authenticate_client_certificate": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"application_entity": schemaApplicationEntity(),
					},
				},
			},
			"store_service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_retrieve_service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDICOMGatewayConfigDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("") // Nothing to do for now
	return diags
}

func resourceDICOMGatewayConfigRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	configURL := d.Get("config_url").(string)
	client, err := config.getDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	// Refresh token so we hopefully have DICOM permissions to proceed without error
	_ = client.TokenRefresh()
	storeConfig, _, err := client.Config.GetStoreService()
	if err != nil {
		return diag.FromErr(err)
	}
	_ = setSCPConfig(*storeConfig, d)

	queryConfig, _, err := client.Config.GetQueryService(nil)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = setQueryConfig(*queryConfig, d)

	return diags
}

func setSCPConfig(scpConfig dicom.SCPConfig, d *schema.ResourceData) error {
	storeService := make(map[string]interface{})
	secure := scpConfig.SecureNetworkConnection.IsSecure
	storeService["is_secure"] = scpConfig.SecureNetworkConnection.IsSecure
	if secure {
		storeService["port"] = scpConfig.SecureNetworkConnection.Port
		if scpConfig.SecureNetworkConnection.AdvancedSettings != nil {
			storeService["pdu_length"] = scpConfig.SecureNetworkConnection.AdvancedSettings.PDULength
			storeService["artim_timeout"] = scpConfig.SecureNetworkConnection.AdvancedSettings.ArtimTimeout
		}
	} else {
		storeService["port"] = scpConfig.UnSecureNetworkConnection.Port
		if scpConfig.UnSecureNetworkConnection.AdvancedSettings != nil {
			storeService["pdu_length"] = scpConfig.UnSecureNetworkConnection.AdvancedSettings.PDULength
			storeService["artim_timeout"] = scpConfig.UnSecureNetworkConnection.AdvancedSettings.ArtimTimeout
		}
	}
	// Add applications
	a := &schema.Set{F: resourceMetricsThresholdHash}
	for _, app := range scpConfig.ApplicationEntities {
		entry := make(map[string]interface{})
		entry["allow_any"] = app.AllowAny
		entry["ae_title"] = app.AeTitle
		entry["organization_id"] = app.OrganizationID
		a.Add(entry)
	}
	storeService["application_entity"] = a

	s := &schema.Set{F: resourceMetricsThresholdHash} // TODO: look at the significance of this
	s.Add(storeService)
	_ = d.Set("store_service", s)
	return nil
}

func setQueryConfig(queryConfig dicom.SCPConfig, d *schema.ResourceData) error {
	queryService := make(map[string]interface{})
	secure := queryConfig.SecureNetworkConnection.IsSecure
	if secure {
		queryService["port"] = queryConfig.SecureNetworkConnection.Port

	} else {
		queryService["port"] = queryConfig.UnSecureNetworkConnection.Port
	}
	// Add applications
	a := &schema.Set{F: resourceMetricsThresholdHash}
	for _, app := range queryConfig.ApplicationEntities {
		entry := make(map[string]interface{})
		entry["allow_any"] = app.AllowAny
		entry["ae_title"] = app.AeTitle
		entry["organization_id"] = app.OrganizationID
		a.Add(entry)
	}
	queryService["application_entity"] = a

	s := &schema.Set{F: resourceMetricsThresholdHash} // TODO: look at the significance of this
	s.Add(queryService)
	_ = d.Set("query_retrieve_service", s)
	return nil
}

func getSCPConfig(d *schema.ResourceData) (*dicom.SCPConfig, error) {
	var scpConfig dicom.SCPConfig

	if v, ok := d.GetOk("store_service"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			scpConfig.Title = mVi["title"].(string)
			scpConfig.Description = mVi["description"].(string)
			isSecure := mVi["is_secure"].(bool)
			port := mVi["port"].(int)
			pduLength := mVi["pdu_length"].(int)
			artimTimeout := mVi["artim_timeout"].(int)
			associationIdleIimeout := mVi["association_idle_timeout"].(int)
			certificateID := mVi["certificate_id"].(string)
			authenticateClientCertificate := mVi["authenticate_client_certificate"].(bool)
			if isSecure {
				if port == 0 {
					port = 105
				}
				scpConfig.SecureNetworkConnection = &dicom.NetworkConnection{
					IsSecure:                      true,
					Port:                          port,
					AuthenticateClientCertificate: authenticateClientCertificate,
					AdvancedSettings: &dicom.AdvancedSettings{
						ArtimTimeout:           artimTimeout,
						AssociationIdleTimeout: associationIdleIimeout,
						PDULength:              pduLength,
					},
				}
				if certificateID != "" {
					scpConfig.SecureNetworkConnection.CertificateInfo = &dicom.CertificateInfo{
						ID: certificateID,
					}
				}
			} else {
				if port == 0 {
					port = 104
				}
				scpConfig.UnSecureNetworkConnection = &dicom.NetworkConnection{
					IsSecure: false,
					Port:     port,
					AdvancedSettings: &dicom.AdvancedSettings{
						ArtimTimeout:           artimTimeout,
						AssociationIdleTimeout: associationIdleIimeout,
						PDULength:              pduLength,
					},
				}
			}
			if as, ok := mVi["application_entity"].(*schema.Set); ok {
				aL := as.List()
				for _, entry := range aL {
					app := entry.(map[string]interface{})
					scpConfig.ApplicationEntities = append(scpConfig.ApplicationEntities, dicom.ApplicationEntity{
						AllowAny:       app["allow_any"].(bool),
						AeTitle:        app["ae_title"].(string),
						OrganizationID: app["organization_id"].(string),
					})
				}
			}
		}
	}

	return &scpConfig, nil
}

func getQueryRetrieveConfig(d *schema.ResourceData) (*dicom.SCPConfig, error) {
	var queryRetrieveConfig dicom.SCPConfig
	if v, ok := d.GetOk("query_retrieve_service"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			queryRetrieveConfig.Title = mVi["title"].(string)
			queryRetrieveConfig.Description = mVi["description"].(string)
			isSecure := mVi["is_secure"].(bool)
			port := mVi["port"].(int)
			pduLength := mVi["pdu_length"].(int)
			artimTimeout := mVi["artim_timeout"].(int)
			associationIdleIimeout := mVi["association_idle_timeout"].(int)
			certificateID := mVi["certificate_id"].(string)
			authenticateClientCertificate := mVi["authenticate_client_certificate"].(bool)
			if isSecure {
				if port == 0 {
					port = 109
				}
				queryRetrieveConfig.SecureNetworkConnection = &dicom.NetworkConnection{
					Port:                          port,
					AuthenticateClientCertificate: authenticateClientCertificate,
					AdvancedSettings: &dicom.AdvancedSettings{
						ArtimTimeout:           artimTimeout,
						AssociationIdleTimeout: associationIdleIimeout,
						PDULength:              pduLength,
					},
				}
				if certificateID != "" {
					queryRetrieveConfig.SecureNetworkConnection.CertificateInfo = &dicom.CertificateInfo{
						ID: certificateID,
					}
				}
			} else {
				if port == 0 {
					port = 108
				}
				queryRetrieveConfig.UnSecureNetworkConnection = &dicom.NetworkConnection{
					Port:     port,
					AdvancedSettings: &dicom.AdvancedSettings{
						ArtimTimeout:           artimTimeout,
						AssociationIdleTimeout: associationIdleIimeout,
						PDULength:              pduLength,
					},
				}
			}
			if as, ok := mVi["application_entity"].(*schema.Set); ok {
				aL := as.List()
				for _, entry := range aL {
					app := entry.(map[string]interface{})
					queryRetrieveConfig.ApplicationEntities = append(queryRetrieveConfig.ApplicationEntities, dicom.ApplicationEntity{
						AllowAny:       app["allow_any"].(bool),
						AeTitle:        app["ae_title"].(string),
						OrganizationID: app["organization_id"].(string),
					})
				}
			}
		}
	}
	return &queryRetrieveConfig, nil
}

func resourceDICOMGatewayConfigCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	configURL := d.Get("config_url").(string)
	client, err := config.getDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	// Refresh token, so we hopefully have DICOM permissions to proceed without error
	_ = client.TokenRefresh()

	scpConfig, err := getSCPConfig(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("getSCPConfig: %w", err))
	}

	queryConfig, err := getQueryRetrieveConfig(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("getQueryRetrieveConfig: %w", err))
	}

	createdSCPConfig, _, err := client.Config.SetStoreService(*scpConfig)
	if err != nil {
		return diag.FromErr(fmt.Errorf("SetStoreService: %w", err))
	}
	_ = d.Set("store_service_id", createdSCPConfig.ID)

	createdQuerySCPConfig, _, err := client.Config.SetQueryService(*queryConfig, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("SetQueryService: %w", err))
	}
	_ = d.Set("query_retrieve_service_id", createdQuerySCPConfig.ID)

	d.SetId(fmt.Sprintf("%s:%s", createdQuerySCPConfig.ID, createdQuerySCPConfig.ID))
	return diags
}
