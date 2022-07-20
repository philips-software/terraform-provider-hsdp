package role_sharing_policy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRoleSharingPolicy() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceRoleSharingPolicyCreate,
		ReadContext:   resourceRoleSharingPolicyRead,
		UpdateContext: resourceRoleSharingPolicyUpdate,
		DeleteContext: resourceRoleSharingPolicyDelete,

		Schema: map[string]*schema.Schema{
			"sharing_policy": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Restricted", "AllowChildren", "Denied",
				}, false),
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"purpose": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_organization_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_organization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRoleSharingPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	roleID := d.Get("role_id").(string)
	targetOrganizationID := d.Get("target_organization_id").(string)
	sharingPolicy := d.Get("sharing_policy").(string)

	var policy *iam.RoleSharingPolicy
	var resp *iam.Response

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var err error
		policy, resp, err = client.Roles.RemoveSharingPolicy(iam.Role{ID: roleID}, iam.RoleSharingPolicy{
			TargetOrganizationID: targetOrganizationID,
			SharingPolicy:        sharingPolicy,
		})
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp == nil {
			return diag.FromErr(fmt.Errorf("response is nil: %v", err))
		}
		return diag.FromErr(err)
	}
	if policy.InternalID != d.Id() {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "internal policy ID mismatch",
			Detail:   fmt.Sprintf("the internal ID of the stored and deleted policy did not match: '%s' -> '%s'", policy.InternalID, d.Id()),
		})
	}
	d.SetId("")
	return diags
}

func resourceRoleSharingPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceRoleSharingPolicyCreate(ctx, d, m)
}

func resourceRoleSharingPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	roleID := d.Get("role_id").(string)
	targetOrganizationID := d.Get("target_organization_id").(string)

	var policies *[]iam.RoleSharingPolicy
	var resp *iam.Response

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		policies, resp, err = client.Roles.ListSharingPolicies(iam.Role{ID: roleID}, &iam.ListSharingPoliciesOptions{
			TargetOrganizationID: &targetOrganizationID,
		})
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp == nil {
			return diag.FromErr(fmt.Errorf("response is nil: %v", err))
		}
		return diag.FromErr(err)
	}
	if len(*policies) == 0 { // No previous tuple found
		d.SetId("")
		return diags
	}
	policy := (*policies)[0]

	_ = d.Set("purpose", policy.Purpose)
	_ = d.Set("source_organization_id", policy.SourceOrganizationID)
	_ = d.Set("sharing_policy", policy.SharingPolicy)
	_ = d.Set("role_name", policy.RoleName)
	_ = d.Set("role_id", policy.RoleID)
	return diags
}

func resourceRoleSharingPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	roleID := d.Get("role_id").(string)
	purpose := d.Get("purpose").(string)
	targetOrganizationID := d.Get("target_organization_id").(string)
	sharingPolicy := d.Get("sharing_policy").(string)

	var policy *iam.RoleSharingPolicy
	var resp *iam.Response

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var err error
		policy, resp, err = client.Roles.ApplySharingPolicy(iam.Role{ID: roleID}, iam.RoleSharingPolicy{
			SharingPolicy:        sharingPolicy,
			TargetOrganizationID: targetOrganizationID,
			Purpose:              purpose,
		})
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp == nil {
			return diag.FromErr(fmt.Errorf("response is nil: %v", err))
		}
		return diag.FromErr(err)
	}

	d.SetId(policy.InternalID)
	return diags
}
