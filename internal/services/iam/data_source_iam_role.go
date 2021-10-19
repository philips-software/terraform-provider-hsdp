package iam

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
	config2 "github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceIAMRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIAMRoleRead,
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

func dataSourceIAMRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config2.Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	orgId := d.Get("managing_organization_id").(string)
	name := d.Get("name").(string)

	roles, resp, err := client.Roles.GetRoles(&iam.GetRolesOptions{
		OrganizationID: &orgId,
		Name:           &name,
	})

	if err != nil {
		if resp != nil && resp.StatusCode != http.StatusOK {
			switch resp.StatusCode {
			case http.StatusNotFound:
				err = fmt.Errorf("role '%s' not found in org '%s'", name, orgId)
			case http.StatusForbidden:
				err = fmt.Errorf("no permission to read roles in org '%s'", orgId)
			}
		}
		return diag.FromErr(err)
	}
	if len(*roles) == 0 || len(*roles) > 1 {
		return diag.FromErr(fmt.Errorf("found %d matching roles", len(*roles)))
	}
	role := (*roles)[0]

	d.SetId(role.ID)
	_ = d.Set("name", role.Name)
	_ = d.Set("description", role.Description)

	return diags
}
