package role_sharing_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-dip-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceIAMRoleSharingPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIAMRoleSharingPoliciesRead,
		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_organization_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sharing_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sharing_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"target_organization_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_organization_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"role_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"purposes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceIAMRoleSharingPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	roleID := d.Get("role_id").(string)
	sharingPolicy := d.Get("sharing_policy").(string)
	targetOrganizationID := d.Get("target_organization_id").(string)

	listOptions := &iam.ListSharingPoliciesOptions{}
	if sharingPolicy != "" {
		listOptions.SharingPolicy = &sharingPolicy
	}
	if targetOrganizationID != "" {
		listOptions.TargetOrganizationID = &targetOrganizationID
	}

	policies, _, err := client.Roles.ListSharingPolicies(iam.Role{ID: roleID}, listOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []string
	var sharingPolicies []string
	var targetOrganizationIDs []string
	var sourceOrganizationIDs []string
	var roleNames []string
	var purposes []string

	for _, policy := range *policies {
		ids = append(ids, policy.InternalID)
		sharingPolicies = append(sharingPolicies, policy.SharingPolicy)
		targetOrganizationIDs = append(targetOrganizationIDs, policy.TargetOrganizationID)
		sourceOrganizationIDs = append(sourceOrganizationIDs, policy.TargetOrganizationID)
		roleNames = append(roleNames, policy.RoleName)
		purposes = append(purposes, policy.Purpose)
	}
	_ = d.Set("ids", ids)
	_ = d.Set("sharing_policies", sharingPolicies)
	_ = d.Set("target_organization_ids", targetOrganizationIDs)
	_ = d.Set("source_organization_ids", sourceOrganizationIDs)
	_ = d.Set("role_names", roleNames)
	_ = d.Set("purposes", purposes)

	d.SetId(roleID + targetOrganizationID)
	return diags
}
