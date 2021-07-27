package hsdp

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCDLDataTypeDefinition() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCDLDataTypeDefinitionRead,
		Schema: map[string]*schema.Schema{
			"cdl_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dtd_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"json_schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceCDLDataTypeDefinitionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	dtdId := d.Get("dtd_id").(string)
	endpoint := d.Get("cdl_endpoint").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	dataTypeDefinition, _, err := client.DataTypeDefinition.GetDataTypeDefinitionByID(dtdId)
	if err != nil {
		return diag.FromErr(err)
	}

	b, err := json.Marshal((*dataTypeDefinition).JsonSchema)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId((*dataTypeDefinition).ID)
	_ = d.Set("json_schema", string(b))
	_ = d.Set("created_by", (*dataTypeDefinition).CreatedBy)
	_ = d.Set("created_on", (*dataTypeDefinition).CreatedOn)
	_ = d.Set("updated_by", (*dataTypeDefinition).UpdatedBy)
	_ = d.Set("updated_on", (*dataTypeDefinition).UpdatedOn)

	return diags
}
