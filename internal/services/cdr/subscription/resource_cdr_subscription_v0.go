package subscription

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Upgrades a CDR Subscription resource from v0 to v1
func patchSubscriptionV0(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		rawState = map[string]interface{}{}
	}
	// New version field in version 1
	rawState["version"] = "stu3"
	return rawState, nil
}

func resourceCDRSubscriptionV0() *schema.Resource {
	return &schema.Resource{
		// This is only used for state migration, so the CRUD
		// callbacks are no longer relevant
		Schema: map[string]*schema.Schema{
			"fhir_store": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delete_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"criteria": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reason": {
				Type:     schema.TypeString,
				Required: true,
			},
			"headers": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"end": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
