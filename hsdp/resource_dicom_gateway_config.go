package hsdp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
					Default:  true,
				},
				"ae_title": {
					Type:     schema.TypeString,
					Required: true,
				},
				"organization_id": {
					Type:     schema.TypeString,
					Required: true,
				},
				"service_timeout": {
					Type:     schema.TypeInt,
					Required: true,
					Default:  0,
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
		UpdateContext: resourceDICOMGatewayConfigUpdate,
		DeleteContext: resourceDICOMGatewayConfigDelete,

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
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_secure": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
							Default:  104,
						},
						"advanced_settings": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pdu_length": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"artim_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"association_idle_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"application_entity": schemaApplicationEntity(),
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"query_retrieve_service": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_secure": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
							Default:  104,
						},
						"application_entity": schemaApplicationEntity(),
					},
				},
			},
		},
	}
}

func resourceDICOMGatewayConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.FromErr(fmt.Errorf("implement me"))
}

func resourceDICOMGatewayConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.FromErr(fmt.Errorf("implement me"))

}

func resourceDICOMGatewayConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.FromErr(fmt.Errorf("implement me"))
}

func resourceDICOMGatewayConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.FromErr(fmt.Errorf("implement me"))
}
