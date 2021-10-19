package cdl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceCDLExportRoute() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCDLExportRouteRead,
		Schema: map[string]*schema.Schema{
			"cdl_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"export_route_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"export_route_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_export": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"destination": {
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

func dataSourceCDLExportRouteRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	exportRouteId := d.Get("export_route_id").(string)
	endpoint := d.Get("cdl_endpoint").(string)

	client, err := c.GetCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	exportRoute, _, err := client.ExportRoute.GetExportRouteByID(exportRouteId)
	if err != nil {
		return diag.FromErr(err)
	} else if exportRoute == nil {
		return diag.FromErr(fmt.Errorf("ExportRoute with ID %s not found", exportRouteId))
	}

	sourceBytes, err := json.Marshal((*exportRoute).Source)
	if err != nil {
		return diag.FromErr(err)
	}

	destinationBytes, err := json.Marshal((*exportRoute).Destination)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId((*exportRoute).ID)
	_ = d.Set("export_route_name", (*exportRoute).ExportRouteName)
	_ = d.Set("description", (*exportRoute).Description)
	_ = d.Set("display_name", (*exportRoute).DisplayName)
	_ = d.Set("source", string(sourceBytes))
	_ = d.Set("auto_export", (*exportRoute).AutoExport)
	_ = d.Set("destination", string(destinationBytes))
	_ = d.Set("created_by", (*exportRoute).CreatedBy)
	_ = d.Set("created_on", (*exportRoute).CreatedOn)
	_ = d.Set("updated_by", (*exportRoute).UpdatedBy)
	_ = d.Set("updated_on", (*exportRoute).UpdatedOn)

	return diags
}
