package iam

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceIAMGroup() *schema.Resource {
	return &schema.Resource{
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
		if resp != nil && resp.StatusCode != http.StatusOK {
			switch resp.StatusCode {
			case http.StatusForbidden:
				err = fmt.Errorf("no permission to read groups in org '%s'", orgId)
			default:
				err = fmt.Errorf("group '%s' not found in org '%s' (code: %d)", name, orgId, resp.StatusCode)
			}
		}
		return diag.FromErr(err)
	}
	if len(*groups) == 0 {
		return diag.FromErr(fmt.Errorf("no group matches the search criteria"))
	}
	group := (*groups)[0]

	d.SetId(group.ID)
	_ = d.Set("name", group.GroupName)
	_ = d.Set("description", group.GroupDescription)

	return diags
}
