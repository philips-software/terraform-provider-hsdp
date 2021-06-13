package hsdp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/pki"
)

func resourcePKITenant() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourcePKITenantCreate,
		ReadContext:   resourcePKITenantRead,
		UpdateContext: resourcePKITenantUpdate,
		DeleteContext: resourcePKITenantDelete,

		Schema: map[string]*schema.Schema{
			"organization_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"space_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"iam_orgs": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"role": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key_bits": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"allow_ip_sans": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"allow_any_name": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"allow_subdomains": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"allowed_domains": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_other_sans": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_serial_numbers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_uri_sans": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"client_flag": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"server_flag": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"enforce_hostnames": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"ca": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true, // Updates are not supported
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"common_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"logical_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePKITenantDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	var err error
	var client *pki.Client

	client, err = config.PKIClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	logicalPath, err := pki.APIEndpoint(d.Id()).LogicalPath()
	if err != nil {
		return diag.FromErr(fmt.Errorf("delete PKI tenant: %w", err))
	}
	tenant, _, err := client.Tenants.Retrieve(logicalPath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("delete PKI tenant retrieve: %w", err))
	}
	tenant.ServiceParameters.LogicalPath = logicalPath
	ok, resp, err := client.Tenants.Offboard(*tenant)
	if err != nil {
		return diag.FromErr(fmt.Errorf("delete PK tenant call: %w", err))
	}
	if !ok {
		diags = append(diags, diag.FromErr(fmt.Errorf("delete returned false, http.StatusCode=%d", resp.StatusCode))...)
	}
	return diags
}

func resourcePKITenantUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	var err error
	var client *pki.Client

	client, err = config.PKIClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	tenant, err := schemaToTenant(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	logicalPath, err := pki.APIEndpoint(d.Id()).LogicalPath()
	if err != nil {
		return diag.FromErr(fmt.Errorf("update PKI tenant: %w", err))
	}
	//logicalPath is already determined
	tenant.ServiceParameters.LogicalPath = logicalPath
	_, _, err = client.Tenants.Update(pki.UpdateTenantRequest{
		ServiceParameters: pki.UpdateServiceParameters{
			LogicalPath: logicalPath,
			IAMOrgs:     tenant.ServiceParameters.IAMOrgs,
			Roles:       tenant.ServiceParameters.Roles,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func schemaToTenant(d *schema.ResourceData, _ interface{}) (*pki.Tenant, error) {
	var tenant pki.Tenant
	tenant.OrganizationName = d.Get("organization_name").(string)
	tenant.SpaceName = d.Get("space_name").(string)
	tenant.ServiceName = "hsdp-pki"
	tenant.PlanName = "standard"
	tenant.ServiceParameters.IAMOrgs = expandStringList(d.Get("iam_orgs").(*schema.Set).List())

	if v, ok := d.GetOk("role"); ok {
		vL := v.(*schema.Set).List()
		role := pki.Role{}
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			role.Name = mVi["name"].(string)
			role.KeyType = mVi["key_type"].(string)
			role.KeyBits = mVi["key_bits"].(int)
			role.ClientFlag = mVi["client_flag"].(bool)
			role.ServerFlag = mVi["server_flag"].(bool)
			role.AllowAnyName = mVi["allow_any_name"].(bool)
			role.AllowIPSans = mVi["allow_ip_sans"].(bool)
			role.AllowAnyName = mVi["allow_any_name"].(bool)
			role.AllowSubdomains = mVi["allow_subdomains"].(bool)
			role.EnforceHostnames = mVi["enforce_hostnames"].(bool)
			role.AllowedDomains = expandStringList(mVi["allowed_domains"].(*schema.Set).List())
			role.AllowedOtherSans = expandStringList(mVi["allowed_other_sans"].(*schema.Set).List())
			if len(role.AllowedOtherSans) == 0 {
				role.AllowedOtherSans = []string{"*"}
			}
			role.AllowedSerialNumbers = expandStringList(mVi["allowed_serial_numbers"].(*schema.Set).List())
			role.AllowedURISans = expandStringList(mVi["allowed_uri_sans"].(*schema.Set).List())
			if len(role.AllowedURISans) == 0 {
				role.AllowedURISans = []string{"*"}
			}
			tenant.ServiceParameters.Roles = append(tenant.ServiceParameters.Roles, role)
		}
	}
	if v, ok := d.GetOk("ca"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			tenant.ServiceParameters.CA.CommonName = mVi["common_name"].(string)
		}
	}
	return &tenant, nil
}

func tenantToSchema(tenant pki.Tenant, logicalPath string, d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	if d == nil || logicalPath == "" {
		return fmt.Errorf("tenantToSchema: logicalPath empty or schema.ResourceData are nil")
	}
	_ = d.Set("api_endpoint", d.Id()) // Same
	_ = d.Set("logical_path", logicalPath)
	_ = d.Set("organization_name", tenant.OrganizationName)
	_ = d.Set("space_name", tenant.SpaceName)
	_ = d.Set("service_name", tenant.ServiceName)
	_ = d.Set("plan_name", tenant.PlanName)
	_ = d.Set("iam_orgs", tenant.ServiceParameters.IAMOrgs)

	if count := len(tenant.ServiceParameters.Roles); count > 0 {
		_, _ = config.Debug("Found %d roles\n", count)
		s := &schema.Set{F: resourceMetricsThresholdHash}
		for _, role := range tenant.ServiceParameters.Roles {
			roleDef := make(map[string]interface{})
			roleDef["name"] = role.Name
			roleDef["key_type"] = role.KeyType
			roleDef["key_bits"] = role.KeyBits
			roleDef["enforce_hostnames"] = role.EnforceHostnames
			roleDef["client_flag"] = role.ClientFlag
			roleDef["server_flag"] = role.ServerFlag
			roleDef["allow_any_name"] = role.AllowAnyName
			roleDef["allow_ip_sans"] = role.AllowIPSans
			roleDef["allow_subdomains"] = role.AllowSubdomains
			roleDef["allowed_domains"] = role.AllowedDomains
			roleDef["allowed_other_sans"] = role.AllowedOtherSans
			roleDef["allowed_serial_numbers"] = role.AllowedSerialNumbers
			roleDef["allowed_uri_sans"] = role.AllowedURISans
			_, _ = config.Debug("Adding role: %s\n", role.Name)
			s.Add(roleDef)
		}
		err := d.Set("role", s)
		if err != nil {
			return err
		}
	} else {
		_, _ = config.Debug("No roles found\n")
	}
	s := &schema.Set{F: resourceMetricsThresholdHash}
	caDef := make(map[string]interface{})
	caDef["common_name"] = tenant.ServiceParameters.CA.CommonName
	s.Add(caDef)
	_ = d.Set("ca", s)
	return nil
}

func resourcePKITenantCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	var err error
	var client *pki.Client

	client, err = config.PKIClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	tenant, err := schemaToTenant(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, _, err := client.Tenants.Onboard(*tenant)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("api_endpoint", resp.APIEndpoint)
	d.SetId(string(resp.APIEndpoint))
	return resourcePKITenantRead(ctx, d, m)
}

func resourcePKITenantRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	var err error
	var client *pki.Client

	client, err = config.PKIClient()
	if err != nil {
		return diag.FromErr(fmt.Errorf("read PKI Tenant client: %w", err))
	}
	defer client.Close()
	logicalPath, err := pki.APIEndpoint(d.Id()).LogicalPath()
	if err != nil {
		return diag.FromErr(fmt.Errorf("read PKI Tenant logical path: %w", err))
	}
	tenant, _, err := client.Tenants.Retrieve(logicalPath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read PKI tenant retrieve: %w", err))
	}
	if err := tenantToSchema(*tenant, logicalPath, d, m); err != nil {
		return diag.FromErr(fmt.Errorf("read PKI Tenant convert to schema: %w", err))
	}
	return diags
}
