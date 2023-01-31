package tenant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

// Upgrades a Tenant from v0 to v1
func patchTenantV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		rawState = map[string]interface{}{}
	}
	if rawState["ca"] == nil {
		rawState["ca"] = []interface{}{}
	}
	for i := range rawState["ca"].([]interface{}) {
		(rawState["ca"].([]interface{})[i]).(map[string]interface{})["ttl"] = "8760h"
	}
	return rawState, nil
}

func ResourcePKITenantV0() *schema.Resource {
	return &schema.Resource{
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
				Elem:     tools.StringSchema(),
			},
			"role": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Required: true,
				Elem:     pkiRoleSchemaV0(),
			},
			"ca": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true, // Updates are not supported
				MaxItems: 1,
				Elem:     pkiCASchemaV0(),
			},
			"triggers": {
				Description: "A map of arbitrary strings that, when changed, will force the resource to be replaced.",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
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

func pkiCASchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"common_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func pkiRoleSchemaV0() *schema.Resource {
	return &schema.Resource{
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
				Elem:     tools.StringSchema(),
			},
			"allowed_other_sans": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     tools.StringSchema(),
			},
			"allowed_serial_numbers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     tools.StringSchema(),
			},
			"allowed_uri_sans": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     tools.StringSchema(),
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
	}
}
