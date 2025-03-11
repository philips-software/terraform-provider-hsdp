package mdm

import (
	"context"
	"fmt"

	"github.com/dip-software/go-dip-api/connect/mdm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func DataSourceConnectMDMStandardService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMStandardServiceRead,
		Schema: map[string]*schema.Schema{
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trusted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     tools.StringSchema(),
			},
			"service_url": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     serviceURLSchema(),
			},
		},
	}

}

func dataSourceConnectMDMStandardServiceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)

	resources, _, err := client.StandardServices.GetStandardServices(&mdm.GetStandardServiceOptions{
		Name: &name,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var resource *mdm.StandardService
	// Match specific name
	for _, r := range *resources {
		if r.Name == name {
			resource = &r
			break
		}
	}
	if len(*resources) == 0 || resource == nil {
		return diag.FromErr(fmt.Errorf("StandardService '%s' not found", name))
	}

	d.SetId(fmt.Sprintf("StandardService/%s", resource.ID))
	_ = d.Set("guid", resource.ID)
	_ = d.Set("name", resource.Name)
	_ = d.Set("description", resource.Description)
	_ = d.Set("tags", resource.Tags)
	if resource.OrganizationGuid != nil && resource.OrganizationGuid.Value != "" {
		value := resource.OrganizationGuid.Value
		if resource.OrganizationGuid.System != "" {
			value = fmt.Sprintf("%s|%s", resource.OrganizationGuid.System, resource.OrganizationGuid.Value)
		}
		_ = d.Set("organization_identifier", value)
	}
	_ = d.Set("trusted", resource.Trusted)
	s := &schema.Set{F: schema.HashResource(serviceURLSchema())}
	for _, serviceURL := range resource.ServiceUrls {
		entry := make(map[string]interface{})
		entry["url"] = serviceURL.URL
		entry["sort_order"] = serviceURL.SortOrder
		if serviceURL.AuthenticationMethodID != nil {
			entry["authentication_method_id"] = serviceURL.AuthenticationMethodID.Reference
		}
		s.Add(entry)
	}
	_ = d.Set("service_url", s)
	return diags
}
