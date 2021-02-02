package hsdp

import (
	"context"
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
		DeleteContext: resourceDPKITenantDelete,

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
			"role": {
				Type:     schema.TypeSet,
				MinItems: 1,
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
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_uri_san": {
							Type:     schema.TypeSet,
							Optional: true,
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
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"common_name": {
							Type:     schema.TypeString,
							Required: true,
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
		},
	}
}

func resourceDPKITenantDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourcePKITenantUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourcePKITenantRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	var err error
	var client *pki.Client

	region := d.Get("region").(string)
	environment := d.Get("environment").(string)
	if region != "" || environment != "" {
		client, err = config.PKIClient(region, environment)
	} else {
		client, err = config.PKIClient()
	}
	if err != nil {
		return diag.FromErr(err)
	}
	var tenant pki.Tenant

	resp, _, err := client.Tenants.Onboard(tenant)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("api_endpoint", resp.APIEndpoint)
	d.SetId(resp.APIEndpoint)
	client.Services.GetPolicyCRL()
	return diags

}

func resourcePKITenantCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
