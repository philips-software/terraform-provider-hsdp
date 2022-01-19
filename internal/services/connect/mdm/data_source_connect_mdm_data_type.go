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

func DataSourceConnectMDMDataType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMDataTypeRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"proposition_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     tools.StringSchema(),
			},
		},
	}

}

func dataSourceConnectMDMDataTypeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	propID := d.Get("proposition_id").(string)
	name := d.Get("name").(string)

	dataTypes, _, err := client.DataTypes.Find(&mdm.GetDataTypeOptions{
		PropositionID: &propID,
		Name:          &name,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if dataTypes == nil || len(*dataTypes) == 0 {
		return diag.FromErr(config.ErrResourceNotFound)
	}
	resource := (*dataTypes)[0]

	d.SetId(fmt.Sprintf("DataType/%s", resource.ID))
	_ = d.Set("guid", resource.ID)
	_ = d.Set("tags", resource.Tags)
	_ = d.Set("description", resource.Description)

	return diags
}
