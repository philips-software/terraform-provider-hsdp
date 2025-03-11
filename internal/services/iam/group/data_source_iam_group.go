package group

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dip-software/go-dip-api/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func DataSourceIAMGroup() *schema.Resource {
	return &schema.Resource{
		Description: descriptions["group"],
		ReadContext: dataSourceIAMGroupRead,
		Schema: map[string]*schema.Schema{
			"managing_organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"services": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"devices": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceIAMGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	orgId := d.Get("managing_organization_id").(string)
	name := d.Get("name").(string)

	groups, resp, err := client.Groups.GetGroups(&iam.GetGroupOptions{
		OrganizationID: &orgId,
		Name:           &name,
	})

	if err != nil {
		if resp != nil && resp.StatusCode() != http.StatusOK {
			switch resp.StatusCode() {
			case http.StatusForbidden:
				err = fmt.Errorf("no permission to read groups in org '%s'", orgId)
			default:
				err = fmt.Errorf("group '%s' not found in org '%s' (code: %d)", name, orgId, resp.StatusCode())
			}
		}
		return diag.FromErr(err)
	}
	if len(*groups) == 0 {
		return diag.FromErr(fmt.Errorf("no group matches the search criteria"))
	}
	group := (*groups)[0]

	_ = d.Set("name", group.GroupName)
	_ = d.Set("description", group.GroupDescription)

	// Extract USER member details
	result, err := getGroupResourcesByMemberType(ctx, client, group.ID, iam.GroupMemberTypeUser)
	if err != nil {
		return diag.FromErr(fmt.Errorf("reading USER members: %w", err))
	}
	_ = d.Set("users", tools.SchemaSetStrings(result))

	// Extract SERVICE member details
	result, err = getGroupResourcesByMemberType(ctx, client, group.ID, iam.GroupMemberTypeService)
	if err != nil {
		return diag.FromErr(fmt.Errorf("reading SERVICE members: %w", err))
	}
	_ = d.Set("services", tools.SchemaSetStrings(result))

	// Extract DEVICE member details
	result, err = getGroupResourcesByMemberType(ctx, client, group.ID, iam.GroupMemberTypeDevice)
	if err != nil {
		return diag.FromErr(fmt.Errorf("reading DEVICE members: %w", err))
	}
	_ = d.Set("devices", tools.SchemaSetStrings(result))

	d.SetId(group.ID)
	return diags
}

func getGroupResourcesByMemberType(ctx context.Context, client *iam.Client, groupID, memberType string) ([]string, error) {
	var resources *iam.SCIMGroup
	var resp *iam.Response
	var err error
	var result []string
	perPage := 100
	err = tools.TryHTTPCall(ctx, 5, func() (*http.Response, error) {
		resources, resp, err = client.Groups.SCIMGetGroupByIDAll(groupID, &iam.SCIMGetGroupOptions{
			IncludeGroupMembersType: &memberType,
			GroupMembersCount:       &perPage,
		})
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return result, err
	}
	if resources == nil {
		return result, fmt.Errorf("unexpected response: %+v", resp)
	}
	for _, u := range resources.ExtensionGroup.GroupMembers.Resources {
		result = append(result, u.ID)
	}
	return result, nil
}
