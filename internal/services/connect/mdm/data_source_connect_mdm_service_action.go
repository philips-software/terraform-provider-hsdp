package mdm

import (
	"context"
	"fmt"

	"github.com/philips-software/go-dip-api/connect/mdm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceConnectMDMServiceAction() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMServiceActionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"standard_service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceConnectMDMServiceActionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)

	resources, _, err := client.ServiceActions.Find(&mdm.GetServiceActionOptions{
		Name: &name,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if resources == nil || len(*resources) == 0 {
		return diag.FromErr(config.ErrResourceNotFound)
	}
	resource := (*resources)[0]

	d.SetId(fmt.Sprintf("ServiceAction/%s", resource.ID))
	_ = d.Set("guid", resource.ID)
	_ = d.Set("type", resource.ResourceType)
	_ = d.Set("description", resource.Description)
	_ = d.Set("standard_service_id", resource.StandardServiceId)
	_ = d.Set("organization_guid", guidValue(resource))

	return diags
}

func guidValue(resource mdm.ServiceAction) string {
	value := ""
	if resource.OrganizationGuid != nil {
		value = resource.OrganizationGuid.Value
		if resource.OrganizationGuid.System != "" {
			value = fmt.Sprintf("%s|%s", resource.OrganizationGuid.System, resource.OrganizationGuid.Value)
		}
	}
	return value
}
