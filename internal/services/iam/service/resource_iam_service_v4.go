package service

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

// Upgrades an IAM Service resource from v4 to v5
func patchIAMServiceV4(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		rawState = map[string]interface{}{}
	}
	// Copy value to new field
	if rawState["self_managed_private_key"] != "" && rawState["expires_on"] != "" {
		rawState["self_managed_expires_on"] = rawState["expires_on"]
	}
	return rawState, nil
}

func ResourceIAMServiceV4() *schema.Resource {
	return &schema.Resource{
		// This is only used for state migration, so the CRUD
		// callbacks are no longer relevant
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressCaseDiffs,
			},
			"description": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"validity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      12,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 600),
			},
			"token_validity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1800,
				ValidateFunc: validation.IntBetween(0, 2592000),
			},
			"self_managed_private_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"self_managed_certificate": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Use 'self_managed_private_key' instead. This will be removed in a future version",
			},
			"private_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_on": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: tools.SuppressWhenGenerated,
			},
			"scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				MinItems: 1, // openid
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"default_scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				MinItems: 1, // openid
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
