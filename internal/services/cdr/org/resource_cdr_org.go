package org

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceCDROrg() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceCDROrgCreate,
		ReadContext:   resourceCDROrgRead,
		UpdateContext: resourceCDROrgUpdate,
		DeleteContext: resourceCDROrgDelete,
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceCDROrgV0().CoreConfigSchema().ImpliedType(),
				Upgrade: patchCDROrgV0,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"fhir_store": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "stu3",
				ForceNew: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"part_of": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"purge_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceCDROrgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	endpoint := d.Get("fhir_store").(string)

	client, err := c.GetFHIRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	version := d.Get("version").(string)

	switch version {
	case "stu3":
		createDiags := stu3Create(ctx, c, client, d, m)
		if len(createDiags) > 0 {
			return createDiags
		}
	case "r4":
		createDiags := r4Create(ctx, c, client, d, m)
		if len(createDiags) > 0 {
			return createDiags
		}
	default:
		return diag.FromErr(fmt.Errorf("unsupported FHIR version '%s'", version))
	}

	return resourceCDROrgRead(ctx, d, m)
}

func resourceCDROrgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics

	endpoint := d.Get("fhir_store").(string)

	client, err := c.GetFHIRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	version := d.Get("version").(string)

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
	return diags
}

func resourceCDROrgUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics

	endpoint := d.Get("fhir_store").(string)

	client, err := c.GetFHIRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	version := d.Get("version").(string)

	switch version {
	case "stu3":
		if updateDiags := stu3Update(ctx, client, d, m); len(updateDiags) > 0 {
			return updateDiags
		}
	case "r4":
		if updateDiags := r4Update(ctx, client, d, m); len(updateDiags) > 0 {
			return updateDiags
		}
	default:
		return diag.FromErr(fmt.Errorf("unsupported FHIR version '%s'", version))
	}
	return diags
}

func resourceCDROrgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics

	endpoint := d.Get("fhir_store").(string)

	client, err := c.GetFHIRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	version := d.Get("version").(string)

	switch version {
	case "stu3":
		if deleteDiags := stu3Delete(ctx, client, d, m); len(deleteDiags) > 0 {
			return deleteDiags
		}
	case "r4":
		if deleteDiags := r4Delete(ctx, client, d, m); len(deleteDiags) > 0 {
			return deleteDiags
		}
	default:
		return diag.FromErr(fmt.Errorf("unsupported FHIR version '%s'", version))
	}
	return diags

}
