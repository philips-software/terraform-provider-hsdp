package service

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

// Upgrades an IAM Service resource from v3 to v4
func patchIAMServiceV3(_ context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		rawState = map[string]interface{}{}
	}
	// New drift_detection field in version 1
	rawState["token_validity"] = 1800
	return rawState, nil
}

func ResourceIAMServiceV3() *schema.Resource {
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
				RequiredWith:     []string{"self_managed_private_key"},
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
