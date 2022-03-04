package mdm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func DataSourceConnectMDMServiceAgent() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMServiceAgentRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"configuration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_subscriber_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_version_supported": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authentication_method_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     tools.StringSchema(),
			},
		},
	}

}

func dataSourceConnectMDMServiceAgentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)

	resources, _, err := client.ServiceAgents.Get(&mdm.GetServiceAgentOptions{
		Name: &name,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if resources == nil || len(*resources) == 0 {
		return diag.FromErr(config.ErrResourceNotFound)
	}
	resource := (*resources)[0]

	d.SetId(fmt.Sprintf("ServiceAgent/%s", resource.ID))
	_ = d.Set("guid", resource.ID)
	_ = d.Set("configuration", resource.Configuration)
	_ = d.Set("description", resource.Description)
	_ = d.Set("authentication_method_ids", resource.AuthenticationMethodIds)
	_ = d.Set("data_subscriber_id", resource.DataSubscriberId)
	_ = d.Set("api_version_supported", resource.APIVersionSupported)

	return diags
}
