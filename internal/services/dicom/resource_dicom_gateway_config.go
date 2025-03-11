package dicom

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/dip-software/go-dip-api/dicom"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func schemaApplicationEntity() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		MaxItems: 100,
		Optional: true,
		Elem:     applicationEntitySchema(),
	}
}

func applicationEntitySchema() *schema.Resource {
	return &schema.Resource{
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
	}
}

func ResourceDICOMGatewayConfig() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceDICOMGatewayConfigCreate,
		ReadContext:   resourceDICOMGatewayConfigRead,
		DeleteContext: resourceDICOMGatewayConfigDelete,
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"config_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"store_service": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     storeServiceSchema(),
			},
			"query_retrieve_service": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     queryRetrieveSchema(),
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

func queryRetrieveSchema() *schema.Resource {
	return &schema.Resource{
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
			"application_entity": schemaApplicationEntity(),
		},
	}
}
func storeServiceSchema() *schema.Resource {
	return &schema.Resource{
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
			// ---Advanced features end
			"application_entity": schemaApplicationEntity(),
		},
	}
}

func resourceDICOMGatewayConfigDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("") // Nothing to do for now
	return diags
}

func resourceDICOMGatewayConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	organizationID := d.Get("organization_id").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	var storeConfig *dicom.BrokenSCPConfig
	var resp *dicom.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		// Refresh token, so we hopefully have DICOM permissions to proceed without error
		_ = client.TokenRefresh()
		storeConfig, resp, err = client.Config.GetStoreService(&dicom.QueryOptions{
			OrganizationID: &organizationID,
		})
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})

	if err != nil {
		return diag.FromErr(err)
	}
	_ = setBrokenSCPConfig(*storeConfig, d)

	var queryConfig *dicom.BrokenSCPConfig
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		_ = client.TokenRefresh()
		queryConfig, resp, err = client.Config.GetQueryRetrieveService(&dicom.QueryOptions{
			OrganizationID: &organizationID,
		})
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})

	if err != nil {
		return diag.FromErr(err)
	}
	_ = setQueryRetrieveConfig(*queryConfig, d)

	return diags
}

func setBrokenSCPConfig(scpConfig dicom.BrokenSCPConfig, d *schema.ResourceData) error {
	storeService := make(map[string]interface{})
	storeService["title"] = scpConfig.Title
	if scpConfig.SecureNetworkConnection != nil {
		storeService["is_secure"] = true
		storeService["port"] = scpConfig.SecureNetworkConnection.Port
		if scpConfig.SecureNetworkConnection.AdvancedSettings != nil {
			storeService["pdu_length"] = scpConfig.SecureNetworkConnection.AdvancedSettings.PDULength
			storeService["artim_timeout"] = scpConfig.SecureNetworkConnection.AdvancedSettings.ArtimTimeout
			storeService["association_idle_timeout"] = scpConfig.SecureNetworkConnection.AdvancedSettings.AssociationIdleTimeout
		}
	}
	if scpConfig.UnSecureNetworkConnection != nil {
		storeService["is_secure"] = false
		storeService["port"] = scpConfig.UnSecureNetworkConnection.Port
		if scpConfig.UnSecureNetworkConnection.AdvancedSettings != nil {
			storeService["pdu_length"] = scpConfig.UnSecureNetworkConnection.AdvancedSettings.PDULength
			storeService["artim_timeout"] = scpConfig.UnSecureNetworkConnection.AdvancedSettings.ArtimTimeout
			storeService["association_idle_timeout"] = scpConfig.UnSecureNetworkConnection.AdvancedSettings.AssociationIdleTimeout
		}
	}
	// Add applications
	a := &schema.Set{F: schema.HashResource(applicationEntitySchema())}
	for _, app := range scpConfig.ApplicationEntities {
		entry := make(map[string]interface{})
		entry["allow_any"] = app.AllowAny
		entry["ae_title"] = app.AeTitle
		entry["organization_id"] = app.OrganizationID
		a.Add(entry)
	}
	storeService["application_entity"] = a

	s := &schema.Set{F: schema.HashResource(storeServiceSchema())}
	s.Add(storeService)
	_ = d.Set("store_service", s)
	return nil
}

func setQueryRetrieveConfig(queryConfig dicom.BrokenSCPConfig, d *schema.ResourceData) error {
	queryService := make(map[string]interface{})
	queryService["title"] = queryConfig.Title
	queryService["description"] = queryConfig.Description
	if queryConfig.SecureNetworkConnection != nil {
		queryService["port"] = queryConfig.SecureNetworkConnection.Port
		queryService["is_secure"] = true
		if queryConfig.SecureNetworkConnection.AdvancedSettings != nil {
			queryService["artim_timeout"] = queryConfig.SecureNetworkConnection.AdvancedSettings.ArtimTimeout
			queryService["association_idle_timeout"] = queryConfig.SecureNetworkConnection.AdvancedSettings.AssociationIdleTimeout
			queryService["pdu_length"] = queryConfig.SecureNetworkConnection.AdvancedSettings.PDULength
		}
	}
	if queryConfig.UnSecureNetworkConnection != nil {
		queryService["port"] = queryConfig.UnSecureNetworkConnection.Port
		queryService["is_secure"] = false
		if queryConfig.UnSecureNetworkConnection.AdvancedSettings != nil {
			queryService["artim_timeout"] = queryConfig.UnSecureNetworkConnection.AdvancedSettings.ArtimTimeout
			queryService["association_idle_timeout"] = queryConfig.UnSecureNetworkConnection.AdvancedSettings.AssociationIdleTimeout
			queryService["pdu_length"] = queryConfig.UnSecureNetworkConnection.AdvancedSettings.PDULength
		}
	}
	// Add applications
	a := &schema.Set{F: schema.HashResource(applicationEntitySchema())}
	for _, app := range queryConfig.ApplicationEntities {
		entry := make(map[string]interface{})
		entry["allow_any"] = app.AllowAny
		entry["ae_title"] = app.AeTitle
		entry["organization_id"] = app.OrganizationID
		a.Add(entry)
	}
	queryService["application_entity"] = a

	s := &schema.Set{F: schema.HashResource(queryRetrieveSchema())}
	s.Add(queryService)
	_ = d.Set("query_retrieve_service", s)
	return nil
}

func getBrokenSCPConfig(d *schema.ResourceData) (*dicom.BrokenSCPConfig, error) {
	var scpConfig dicom.BrokenSCPConfig

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
			if isSecure {
				if port == 0 {
					port = 105
				}
				scpConfig.SecureNetworkConnection = &dicom.BrokenNetworkConnection{
					Port: port,
					AdvancedSettings: &dicom.BrokenAdvancedSettings{
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
				scpConfig.UnSecureNetworkConnection = &dicom.BrokenNetworkConnection{
					Port: port,
					AdvancedSettings: &dicom.BrokenAdvancedSettings{
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

func getQueryRetrieveConfig(d *schema.ResourceData) (*dicom.BrokenSCPConfig, error) {
	var queryRetrieveConfig dicom.BrokenSCPConfig
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
			if isSecure {
				if port == 0 {
					port = 109
				}
				queryRetrieveConfig.SecureNetworkConnection = &dicom.BrokenNetworkConnection{
					Port: port,
					AdvancedSettings: &dicom.BrokenAdvancedSettings{
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
				queryRetrieveConfig.UnSecureNetworkConnection = &dicom.BrokenNetworkConnection{
					Port: port,
					AdvancedSettings: &dicom.BrokenAdvancedSettings{
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

func resourceDICOMGatewayConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	organizationID := d.Get("organization_id").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	// Refresh token, so we hopefully have DICOM permissions to proceed without error
	_ = client.TokenRefresh()

	scpConfig, err := getBrokenSCPConfig(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("getSCPConfig: %w", err))
	}

	queryConfig, err := getQueryRetrieveConfig(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("getQueryRetrieveConfig: %w", err))
	}

	createdSCPConfig, _, err := client.Config.SetStoreService(*scpConfig, &dicom.QueryOptions{
		OrganizationID: &organizationID,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("SetStoreService: %w", err))
	}
	_ = d.Set("store_service_id", createdSCPConfig.ID)

	createdQuerySCPConfig, _, err := client.Config.SetQueryRetrieveService(*queryConfig, &dicom.QueryOptions{
		OrganizationID: &organizationID,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("SetMoveService: %w", err))
	}
	_ = d.Set("query_retrieve_service_id", createdQuerySCPConfig.ID)

	generatedID := fmt.Sprintf("%x", md5.Sum([]byte(configURL)))
	d.SetId(generatedID)
	return resourceDICOMGatewayConfigRead(ctx, d, m)
}
