package hsdp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func dataSourceIAMUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIAMUsersRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filter": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				MinItems: 0,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_verified": {
							Type:     schema.TypeBool,
							Default:  true,
							Optional: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
					},
				},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"logins": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"email_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

type iamUserFilter struct {
	EmailVerified bool
	Disabled      bool
}

func collectFilter(d *schema.ResourceData) (*iamUserFilter, diag.Diagnostics) {
	var diags diag.Diagnostics
	var filter *iamUserFilter
	if v, ok := d.GetOk("filter"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			filter = &iamUserFilter{
				EmailVerified: mVi["email_verified"].(bool),
				Disabled:      mVi["disabled"].(bool),
			}
		}
	}
	return filter, diags
}

func dataSourceIAMUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	orgID := d.Get("organization_id").(string)
	filter, filterDiags := collectFilter(d)
	if len(filterDiags) > 0 {
		return filterDiags
	}
	profileType := "all"

	userList, _, err := client.Users.GetAllUsers(&iam.GetUserOptions{
		OrganizationID: &orgID,
		ProfileType:    &profileType,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []string
	var logins []string
	var emailAddresses []string

	for _, guid := range userList {
		user, _, err := client.Users.GetUserByID(guid)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("skipped user '%s' due to error", guid),
				Detail:   err.Error(),
			})
			continue
		}
		if filter != nil && user.AccountStatus.Disabled != filter.Disabled {
			continue
		}
		if filter != nil && user.AccountStatus.EmailVerified != filter.EmailVerified {
			continue
		}
		ids = append(ids, guid)
		logins = append(logins, user.LoginID)
		emailAddresses = append(emailAddresses, user.EmailAddress)
	}
	_ = d.Set("ids", ids)
	_ = d.Set("logins", logins)
	_ = d.Set("email_addresses", emailAddresses)
	d.SetId(orgID)
	return diags
}
