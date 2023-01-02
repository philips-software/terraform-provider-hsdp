package application

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

// Upgrades an IAM Group resource from v0 to v1
func patchIAMApplicationV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		rawState = map[string]interface{}{}
	}
	return rawState, nil
}

func ResourceIAMApplicationV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tools.ValidateUpperString,
				ForceNew:     true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"proposition_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"global_reference_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressWhenGenerated,
			},
		},
	}
}
