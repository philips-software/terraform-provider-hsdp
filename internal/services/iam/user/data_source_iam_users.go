package user

import (
	"context"
	"fmt"

	"github.com/dip-software/go-dip-api/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceIAMUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIAMUsersRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_verified": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
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

func dataSourceIAMUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	orgID := d.Get("organization_id").(string)
	var emailVerified *bool
	var disabled *bool
	if val, ok := d.GetOk("disabled"); ok {
		b := val.(bool)
		disabled = &b
	}
	if val, ok := d.GetOk("email_verified"); ok {
		b := val.(bool)
		emailVerified = &b
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
		if emailVerified != nil && *emailVerified != user.AccountStatus.EmailVerified {
			continue
		}
		if disabled != nil && *disabled != user.AccountStatus.Disabled {
			continue
		}
		// All criteria match, so add user
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
