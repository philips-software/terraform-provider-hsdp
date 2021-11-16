package mdm

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceConnectMDMStorageClasses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMStorageClassesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceConnectMDMStorageClassesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resources, _, err := client.StorageClasses.GetStorageClasses(nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []string
	var names []string
	var descriptions []string

	for _, r := range *resources {
		ids = append(ids, r.ID)
		names = append(names, r.Name)
		descriptions = append(descriptions, r.Description)
	}
	_ = d.Set("ids", ids)
	_ = d.Set("names", names)
	_ = d.Set("descriptions", descriptions)
	result, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(result)
	return diags
}
