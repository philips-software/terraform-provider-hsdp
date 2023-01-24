package org

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceCDROrg() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCDROrgRead,

		Schema: map[string]*schema.Schema{
			"fhir_store": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "stu3",
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"part_of": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCDROrgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics

	endpoint := d.Get("fhir_store").(string)

	client, err := c.GetFHIRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	version := d.Get("version").(string)

	id := d.Get("org_id").(string)

	switch version {
	case "stu3":
		if readDiags := stu3Read(ctx, client, d); len(readDiags) > 0 {
			return readDiags
		}
	case "r4":
		if readDiags := r4Read(ctx, client, d); len(readDiags) > 0 {
			return readDiags
		}
	default:
		return diag.FromErr(fmt.Errorf("unsupported FHIR version '%s'", version))
	}
	if len(diags) == 0 {
		d.SetId(id)
	}
	return diags
}
