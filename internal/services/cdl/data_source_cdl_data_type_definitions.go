package cdl

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/cdl"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceCDLDataTypeDefinitions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCDLDataTypeDefinitionsRead,
		Schema: map[string]*schema.Schema{
			"cdl_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceCDLDataTypeDefinitionsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	endpoint := d.Get("cdl_endpoint").(string)

	client, err := c.GetCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	dataTypeDefinitions, _, err := client.DataTypeDefinition.GetDataTypeDefinitions(&cdl.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(endpoint + "datatypedefinitions")

	var dataTypeDefNames []string
	var dataTypeDefIds []string

	for _, dtd := range dataTypeDefinitions {
		dataTypeDefNames = append(dataTypeDefNames, dtd.Name)
		dataTypeDefIds = append(dataTypeDefIds, dtd.ID)
	}
	_ = d.Set("names", dataTypeDefNames)
	_ = d.Set("ids", dataTypeDefIds)

	return diags
}
