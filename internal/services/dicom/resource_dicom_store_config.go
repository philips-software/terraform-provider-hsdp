package dicom

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/dicom"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceDICOMStoreConfig() *schema.Resource {
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
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cdr_service_account": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     cdrSettingsSchema(),
			},
			"fhir_store": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     fhirStoreSettingsSchema(),
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

func fhirStoreSettingsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"mpi_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func cdrSettingsSchema() *schema.Resource {
	return &schema.Resource{
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
	}
}

func resourceDICOMStoreConfigDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}

func resourceDICOMStoreConfigUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var resp *dicom.Response
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	orgID := d.Get("organization_id").(string)
	client, err := c.GetDICOMConfigClient(configURL)
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
			var configured *dicom.CDRServiceAccount
			operation := func() error {
				configured, resp, err = client.Config.SetCDRServiceAccount(cdrService, &dicom.QueryOptions{
					OrganizationID: &orgID,
				})
				return checkForPermissionErrors(client, resp, err)
			}
			err := backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
			if err != nil {
				return diag.FromErr(err)
			}
			cdrSettings := make(map[string]interface{})
			cdrSettings["service_id"] = configured.ServiceID
			cdrSettings["private_key"] = configured.PrivateKey
			cdrSettings["id"] = configured.ID
			s := &schema.Set{F: schema.HashResource(cdrSettingsSchema())}
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
			var configured *dicom.FHIRStore
			operation := func() error {
				configured, resp, err = client.Config.SetFHIRStore(fhirStore, &dicom.QueryOptions{
					OrganizationID: &orgID,
				})
				return checkForPermissionErrors(client, resp, err)
			}
			err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
			if err != nil {
				return diag.FromErr(err)
			}
			fhirSettings := make(map[string]interface{})
			fhirSettings["mpi_endpoint"] = configured.MPIEndpoint
			fhirSettings["id"] = configured.ID
			s := &schema.Set{F: schema.HashResource(fhirStoreSettingsSchema())}
			s.Add(fhirSettings)
			_ = d.Set("fhir_store", s)
		}
	}
	return diags
}

func resourceDICOMStoreConfigRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	orgID := d.Get("organization_id").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	// CDR
	var configured *dicom.CDRServiceAccount
	operation := func() error {
		var resp *dicom.Response
		configured, resp, err = client.Config.GetCDRServiceAccount(&dicom.QueryOptions{
			OrganizationID: &orgID,
		})
		return checkForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewConstantBackOff(2*time.Second), 5))
	if err == nil && configured != nil {
		cdrSettings := make(map[string]interface{})
		cdrSettings["service_id"] = configured.ServiceID
		cdrSettings["private_key"] = configured.PrivateKey
		cdrSettings["id"] = configured.ID
		s := &schema.Set{F: schema.HashResource(cdrSettingsSchema())}
		s.Add(cdrSettings)
		_ = d.Set("cdr_service_account", s)
	}
	// FHIR
	fhirConfigured, _, err := client.Config.GetFHIRStore(&dicom.QueryOptions{
		OrganizationID: &orgID,
	})
	if err == nil && fhirConfigured != nil {
		fhirSettings := make(map[string]interface{})
		fhirSettings["mpi_endpoint"] = fhirConfigured.MPIEndpoint
		fhirSettings["id"] = fhirConfigured.ID
		s := &schema.Set{F: schema.HashResource(fhirStoreSettingsSchema())}
		s.Add(fhirSettings)
		_ = d.Set("fhir_store", s)
	}
	return diags
}

func resourceDICOMStoreConfigCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	orgID := d.Get("organization_id").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	// Refresh token, so we hopefully have DICOM permissions to proceed without error
	_, _ = c.Debug("resourceDICOMStoreConfigCreate: TokenRefresh()\n")
	_ = client.TokenRefresh()
	_, _ = c.Debug("resourceDICOMStoreConfigCreate: refreshed!\n")

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
		var configured *dicom.CDRServiceAccount
		operation := func() error {
			var resp *dicom.Response
			_, _ = c.Debug("resourceDICOMStoreConfigCreate: cdr_service_account operation run\n")
			configured, resp, err = client.Config.SetCDRServiceAccount(cdrService, &dicom.QueryOptions{
				OrganizationID: &orgID,
			})
			return checkForPermissionErrors(client, resp, err)
		}
		err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
		if err != nil {
			return diag.FromErr(err)
		}
		cdrSettings := make(map[string]interface{})
		cdrSettings["service_id"] = configured.ServiceID
		cdrSettings["private_key"] = configured.PrivateKey
		cdrSettings["id"] = configured.ID
		s := &schema.Set{F: schema.HashResource(cdrSettingsSchema())}
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
		var configured *dicom.FHIRStore
		operation := func() error {
			var resp *dicom.Response
			_, _ = c.Debug("resourceDICOMStoreConfigCreate: fhir_store operation run\n")
			configured, _, err = client.Config.SetFHIRStore(fhirStore, &dicom.QueryOptions{
				OrganizationID: &orgID,
			})
			return checkForPermissionErrors(client, resp, err)
		}
		err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
		if err != nil {
			return diag.FromErr(err)
		}
		fhirSettings := make(map[string]interface{})
		fhirSettings["mpi_endpoint"] = configured.MPIEndpoint
		fhirSettings["id"] = configured.ID
		s := &schema.Set{F: schema.HashResource(fhirStoreSettingsSchema())}
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

func checkForPermissionErrors(client *dicom.Client, resp *dicom.Response, err error) error {
	if resp == nil {
		if err == nil {
			return backoff.Permanent(fmt.Errorf("response is 'nil'"))
		}
		return backoff.Permanent(err)
	}
	if resp.StatusCode > 500 {
		return err
	}
	if resp.StatusCode == http.StatusForbidden {
		_ = client.TokenRefresh()
		return err
	}
	return backoff.Permanent(err)
}
