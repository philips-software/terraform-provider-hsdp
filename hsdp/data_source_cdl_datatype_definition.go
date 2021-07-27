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
			"jsonschema": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"createdby": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"createdon": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updatedby": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updatedon": {
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
	_ = d.Set("jsonschema", string(b))
	_ = d.Set("createdby", (*dataTypeDefinition).CreatedBy)
	_ = d.Set("createdon", (*dataTypeDefinition).CreatedOn)
	_ = d.Set("updatedby", (*dataTypeDefinition).UpdatedBy)
	_ = d.Set("updatedon", (*dataTypeDefinition).UpdatedOn)

	return diags
}
