package group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

// Upgrades an IAM Group resource from v1 to v2
func patchIAMGroupV1(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		rawState = map[string]interface{}{}
	}
	return rawState, nil
}

func ResourceIAMGroupV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressCaseDiffs,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tools.SuppressWhenGenerated,
			},
			"managing_organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"roles": {
				Type:     schema.TypeSet,
				MaxItems: 1000,
				Required: true,
				Elem:     tools.StringSchema(),
			},
			"users": {
				Type:     schema.TypeSet,
				MaxItems: 2000,
				Optional: true,
				Elem:     tools.StringSchema(),
			},
			"services": {
				Type:     schema.TypeSet,
				MaxItems: 2000,
				Optional: true,
				Elem:     tools.StringSchema(),
			},
			"drift_detection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}
