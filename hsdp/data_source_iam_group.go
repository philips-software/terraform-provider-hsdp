package hsdp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func dataSourceIAMGroup() *schema.Resource {
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
	config := meta.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	orgId := d.Get("managing_organization_id").(string)
	name := d.Get("name").(string)

	group, resp, err := client.Groups.GetGroup(&iam.GetGroupOptions{
		OrganizationID: &orgId,
		Name:           &name,
	})

	if err != nil {
		if resp != nil && resp.StatusCode != http.StatusOK {
			switch resp.StatusCode {
			case http.StatusNotFound:
				err = fmt.Errorf("group '%s' not found in org '%s'", name, orgId)
			case http.StatusForbidden:
				err = fmt.Errorf("no permission to read groups in org '%s'", orgId)
			}
		}
		return diag.FromErr(err)
	}

	d.SetId(group.ID)
	_ = d.Set("name", group.Name)
	_ = d.Set("description", group.Description)

	return diags
}
