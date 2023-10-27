package group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

// Upgrades an IAM Group resource from v3 to v4
func patchIAMGroupV3(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		rawState = map[string]interface{}{}
	}
	return rawState, nil
}

func ResourceIAMGroupV3() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressCaseDiffs,
				Description:      "The group name.",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tools.SuppressWhenGenerated,
				Description:      "The group description.",
			},
			"managing_organization": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The managing organization ID.",
			},
			"roles": {
				Type:        schema.TypeSet,
				MaxItems:    1000,
				Required:    true,
				Elem:        tools.StringSchema(),
				Description: "The list of role IDS to assign to this group.",
			},
			"users": {
				Type:        schema.TypeSet,
				MaxItems:    2000,
				Optional:    true,
				Elem:        tools.StringSchema(),
				Description: "The list of user IDs to include in this group. The provider only manages this list of users. Existing users added by others means to the group by the provider. It is not practical to manage hundreds or thousands of users this way of course.",
			},
			"services": {
				Type:        schema.TypeSet,
				MaxItems:    2000,
				Optional:    true,
				Elem:        tools.StringSchema(),
				Description: "The list of service identity IDs to include in this group.",
			},
			"devices": {
				Type:        schema.TypeSet,
				MaxItems:    2000,
				Optional:    true,
				Elem:        tools.StringSchema(),
				Description: "The list of IAM device identity IDs to include in this group.",
			},
			"drift_detection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When enabled, the provider will perform additional API calls to determine if any changes were made outside of Terraform to user and service assignments of this group.",
			},
			"iam_device_bug_workaround": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Deprecated:  "This workaround is no longer required and will be removed in the near future.",
				Description: "Deprecated, do not use.",
			},
		},
	}
}
