package hsdp

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSTLDeviceConfig() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceSTLDeviceConfigCreate,
		ReadContext:   resourceSTLDeviceConfigRead,
		UpdateContext: resourceSTLDeviceConfigUpdate,
		DeleteContext: resourceSTLDeviceConfigDelete,

		Schema: map[string]*schema.Schema{
			"serial_number": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"firewall_exceptions": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tcp": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"udp": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"logging": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"raw_config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hsdp_product_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hsdp_shared_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hsdp_secret_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hsdp_logging_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hsdp_custom_field": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"cert": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 100,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"private_key_pem": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cert_pem": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceSTLDeviceConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceSTLDeviceConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceSTLDeviceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceSTLDeviceConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
