package hsdp

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/dicom"
)

func resourceDICOMStoreConfig() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceDICOMStoreConfigCreate,
		ReadContext:   resourceDICOMStoreConfigRead,
		UpdateContext: resourceDICOMStoreConfigUpdate,
		DeleteContext: resourceDICOMStoreConfigDelete,

		Schema: map[string]*schema.Schema{
			"config_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cdr_service_account": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"private_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"fhir_store": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mpi_endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"qido_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stow_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"wado_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_management_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDICOMStoreConfigDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}

func resourceDICOMStoreConfigUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	configURL := d.Get("config_url").(string)
	client, err := config.getDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	if d.HasChange("cdr_service_account") {
		cdrService := dicom.CDRServiceAccount{}
		if v, ok := d.GetOk("cdr_service_account"); ok {
			vL := v.(*schema.Set).List()
			for _, vi := range vL {
				mVi := vi.(map[string]interface{})
				cdrService.ServiceID = mVi["service_id"].(string)
				cdrService.PrivateKey = mVi["private_key"].(string)
			}

			if !cdrService.Valid() {
				return diag.FromErr(fmt.Errorf("cdr_service_account is not valid"))
			}
			configured, _, err := client.Config.SetCDRServiceAccount(cdrService)
			if err != nil {
				return diag.FromErr(err)
			}
			cdrSettings := make(map[string]interface{})
			cdrSettings["service_id"] = configured.ServiceID
			cdrSettings["private_key"] = configured.PrivateKey
			cdrSettings["id"] = configured.ID
			s := &schema.Set{F: resourceMetricsThresholdHash}
			s.Add(cdrSettings)
			_ = d.Set("cdr_service_account", s)
		}
	}
	if d.HasChange("fhir_store") {
		fhirStore := dicom.FHIRStore{}
		if v, ok := d.GetOk("fhir_store"); ok {
			vL := v.(*schema.Set).List()
			for _, vi := range vL {
				mVi := vi.(map[string]interface{})
				fhirStore.MPIEndpoint = mVi["mpi_endpoint"].(string)
			}
			if !fhirStore.Valid() {
				return diag.FromErr(fmt.Errorf("fhir_store is not valid"))
			}
			configured, _, err := client.Config.SetFHIRStore(fhirStore)
			if err != nil {
				return diag.FromErr(err)
			}
			fhirSettings := make(map[string]interface{})
			fhirSettings["mpi_endpoint"] = configured.MPIEndpoint
			fhirSettings["id"] = configured.ID
			s := &schema.Set{F: resourceMetricsThresholdHash}
			s.Add(fhirSettings)
			_ = d.Set("fhir_store", s)
		}
	}
	return diags
}

func resourceDICOMStoreConfigRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	configURL := d.Get("config_url").(string)
	client, err := config.getDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	// CDR
	configured, _, err := client.Config.GetCDRServiceAccount()
	if err != nil {
		return diag.FromErr(err)
	}
	cdrSettings := make(map[string]interface{})
	cdrSettings["service_id"] = configured.ServiceID
	cdrSettings["private_key"] = configured.PrivateKey
	cdrSettings["id"] = configured.ID
	s := &schema.Set{F: resourceMetricsThresholdHash}
	s.Add(cdrSettings)
	_ = d.Set("cdr_service_account", s)

	// FHIR
	fhirConfigured, _, err := client.Config.GetFHIRStore()
	if err != nil {
		return diag.FromErr(err)
	}
	fhirSettings := make(map[string]interface{})
	fhirSettings["mpi_endpoint"] = fhirConfigured.MPIEndpoint
	fhirSettings["id"] = fhirConfigured.ID
	s = &schema.Set{F: resourceMetricsThresholdHash}
	s.Add(fhirSettings)
	_ = d.Set("fhir_store", s)

	return diags
}

func resourceDICOMStoreConfigCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	configURL := d.Get("config_url").(string)
	client, err := config.getDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	// Set up CDR service account
	cdrService := dicom.CDRServiceAccount{}
	if v, ok := d.GetOk("cdr_service_account"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			cdrService.ServiceID = mVi["service_id"].(string)
			cdrService.PrivateKey = mVi["private_key"].(string)
		}

		if !cdrService.Valid() {
			return diag.FromErr(fmt.Errorf("cdr_service_account is not valid"))
		}
		configured, _, err := client.Config.SetCDRServiceAccount(cdrService)
		if err != nil {
			return diag.FromErr(err)
		}
		cdrSettings := make(map[string]interface{})
		cdrSettings["service_id"] = configured.ServiceID
		cdrSettings["private_key"] = configured.PrivateKey
		cdrSettings["id"] = configured.ID
		s := &schema.Set{F: resourceMetricsThresholdHash}
		s.Add(cdrSettings)
		_ = d.Set("cdr_service_account", s)
	}
	// Set up FHIR store
	fhirStore := dicom.FHIRStore{}
	if v, ok := d.GetOk("fhir_store"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			fhirStore.MPIEndpoint = mVi["mpi_endpoint"].(string)
		}
		if !fhirStore.Valid() {
			return diag.FromErr(fmt.Errorf("fhir_store is not valid"))
		}
		configured, _, err := client.Config.SetFHIRStore(fhirStore)
		if err != nil {
			return diag.FromErr(err)
		}
		fhirSettings := make(map[string]interface{})
		fhirSettings["mpi_endpoint"] = configured.MPIEndpoint
		fhirSettings["id"] = configured.ID
		s := &schema.Set{F: resourceMetricsThresholdHash}
		s.Add(fhirSettings)
		_ = d.Set("fhir_store", s)
	}

	// Set URLs
	_ = d.Set("qido_url", client.GetQIDOURL())
	_ = d.Set("wado_url", client.GetWADOURL())
	_ = d.Set("stow_url", client.GetSTOWURL())

	generatedID := fmt.Sprintf("%x", md5.Sum([]byte(configURL)))
	d.SetId(generatedID)
	return diags
}
