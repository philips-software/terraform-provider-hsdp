package mdm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-dip-api/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceConnectMDMStorageClass() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMStorageClassRead,
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
		},
	}

}

func dataSourceConnectMDMStorageClassRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(fmt.Errorf("get MDMClient error: %w", err))
	}

	name := d.Get("name").(string)

	resources, _, err := client.StorageClasses.GetStorageClasses(&mdm.GetStorageClassOptions{
		Name: &name,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("get StorageClasses error: %w", err))
	}

	if len(*resources) == 0 {
		return diag.FromErr(fmt.Errorf("resource '%s' not found", name))
	}
	resource := (*resources)[0]
	d.SetId(fmt.Sprintf("StorageClass/%s", resource.ID))
	_ = d.Set("guid", resource.ID)
	_ = d.Set("name", resource.Name)
	_ = d.Set("description", resource.Description)
	return diags
}
