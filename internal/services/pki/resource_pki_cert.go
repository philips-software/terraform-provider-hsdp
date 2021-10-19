package pki

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/pki"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourcePKICert() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourcePKICertCreate,
		ReadContext:   resourcePKICertRead,
		DeleteContext: resourcePKICertDelete,

		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"common_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alt_name": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Use alt_names, this field is ignored and will be removed in the next major release",
			},
			"alt_names": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_sans": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				Elem:     stringSchema(),
			},
			"uri_sans": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				Elem:     stringSchema(),
			},
			"other_sans": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     stringSchema(),
			},
			"ttl": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"exclude_cn_from_sans": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"cert_pem": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ca_chain_pem": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key_pem": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"expiration": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"issuing_ca_pem": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePKICertCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	var err error
	var client *pki.Client

	client, err = c.PKIClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	tenantID := d.Get("tenant_id").(string)
	logicalPath, err := pki.APIEndpoint(tenantID).LogicalPath()
	if err != nil {
		return diag.FromErr(fmt.Errorf("create PKI cert logicalPath: %w", err))
	}
	tenant, _, err := client.Tenants.Retrieve(logicalPath)
	if err != nil {
		return diag.FromErr(err)
	}
	roleName := d.Get("role").(string)
	ttl := d.Get("ttl").(string)
	ipSANS := tools.ExpandStringList(d.Get("ip_sans").(*schema.Set).List())
	uriSANS := tools.ExpandStringList(d.Get("uri_sans").(*schema.Set).List())
	otherSANS := tools.ExpandStringList(d.Get("other_sans").(*schema.Set).List())
	commonName := d.Get("common_name").(string)
	altNames := d.Get("alt_names").(string)
	excludeCNFromSANS := d.Get("exclude_cn_from_sans").(bool)
	role, ok := tenant.GetRoleOk(roleName)
	if !ok {
		return diag.FromErr(fmt.Errorf("role '%s' not found or invalid", roleName))
	}
	certRequest := pki.CertificateRequest{
		CommonName:        commonName,
		AltNames:          altNames,
		IPSANS:            strings.Join(ipSANS, ","),
		URISANS:           strings.Join(uriSANS, ","),
		OtherSANS:         strings.Join(otherSANS, ","),
		TTL:               ttl,
		ExcludeCNFromSANS: &excludeCNFromSANS,
		PrivateKeyFormat:  "pem",
		Format:            "pem",
	}
	cert, resp, err := client.Services.IssueCertificate(logicalPath, role.Name, certRequest)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusForbidden {
			return diag.FromErr(fmt.Errorf("you might be missing the 'PKI_CERT.ISSUE' permission for the tenant org: %w", err))
		}
		return diag.FromErr(fmt.Errorf("issue PKI cert: %w", err))
	}
	d.SetId(cert.Data.SerialNumber)
	err = certToSchema(cert, d, m)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	return diags
}

func certToSchema(cert *pki.IssueResponse, d *schema.ResourceData, _ interface{}) error {
	if cert.Data.PrivateKey != "" {
		_ = d.Set("private_key_pem", cert.Data.PrivateKey)
	}
	if len(cert.Data.SerialNumber) > 0 {
		_ = d.Set("serial_number", cert.Data.SerialNumber)
	}
	if len(cert.Data.IssuingCa) > 0 {
		_ = d.Set("issuing_ca_pem", cert.Data.IssuingCa)
	}
	if len(cert.Data.Certificate) > 0 {
		_ = d.Set("cert_pem", cert.Data.Certificate)
	}
	if cert.Data.Expiration > 0 {
		_ = d.Set("expiration", cert.Data.Expiration)
	}
	if len(cert.Data.CaChain) > 0 {
		_ = d.Set("ca_chain_pem", strings.Join(cert.Data.CaChain, "\n"))
	}
	return nil
}

func resourcePKICertRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	var err error
	var client *pki.Client

	client, err = c.PKIClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	tenantID := d.Get("tenant_id").(string)
	logicalPath, err := pki.APIEndpoint(tenantID).LogicalPath()
	if err != nil {
		return diag.FromErr(fmt.Errorf("read PKI cert logicalPath: %w", err))
	}
	cert, resp, err := client.Services.GetCertificateBySerial(logicalPath, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound { // Expired, pruned
			d.SetId("")
			return diags
		}
		return diag.FromErr(fmt.Errorf("read PKI cert: %w", err))
	}
	err = certToSchema(cert, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourcePKICertDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	var err error
	var client *pki.Client

	client, err = c.PKIClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	tenantID := d.Get("tenant_id").(string)
	logicalPath, err := pki.APIEndpoint(tenantID).LogicalPath()
	if err != nil {
		return diag.FromErr(fmt.Errorf("delete PKI cert logicalPath: %w", err))
	}
	revoke, _, err := client.Services.RevokeCertificateBySerial(logicalPath, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("delete PKI cert: %w", err))
	}
	if revoke.Data.RevocationTime > 0 {
		d.SetId("")
	}
	return diags
}
