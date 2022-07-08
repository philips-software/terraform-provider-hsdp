package client

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

// Upgrades an IAM Group resource from v0 to v1
func patchIAMClientV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		rawState = map[string]interface{}{}
	}
	return rawState, nil
}

func ResourceIAMClientV0() *schema.Resource {
	return &schema.Resource{
		// This is only used for state migration, so the CRUD
		// callbacks are no longer relevant
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"client_id": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				DiffSuppressFunc: tools.SuppressCaseDiffs,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
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
			"global_reference_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"redirection_uris": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"response_types": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"default_scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"consent_implied": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"access_token_lifetime": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1800,
			},
			"refresh_token_lifetime": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2592000,
			},
			"id_token_lifetime": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3600,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}
