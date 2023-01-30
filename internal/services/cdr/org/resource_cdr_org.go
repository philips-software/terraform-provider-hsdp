package org

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceCDROrg() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: importFHIROrgContext,
		},

		Description:   "Manage CDR Organizations.",
		CreateContext: resourceCDROrgCreate,
		ReadContext:   resourceCDROrgRead,
		UpdateContext: resourceCDROrgUpdate,
		DeleteContext: resourceCDROrgDelete,
		CustomizeDiff: resourceCDROrgDiff(),
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
				Description: "The CDR FHIR store tenant URL.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"version": {
				Description: "The FHIR version to use. Options: [ 'stu3' | 'r4' ].",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "stu3",
			},
			"org_id": {
				Description: "The organization ID (GUID) under which to onboard. Typically the same as IAM Organization ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of the FHIR Organization.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"part_of": {
				Description: "The parent Organization ID (GUID) this Org is part of.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"purge_delete": {
				Description: "If set to true, when the resource is destroyed the provider will purge all FHIR resources associated with the Organization. The ORGANIZATION.PURGE IAM permission is required for this to work.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceCDROrgDiff() schema.CustomizeDiffFunc {
	return customdiff.All(
		customdiff.ValidateChange("version", func(ctx context.Context, old, new, m interface{}) error {
			o := old.(string)
			n := new.(string)
			if o == "" { // New
				return nil
			}
			if o != n {
				return fmt.Errorf("changing the FHIR version from '%s' to '%s' is not supported", o, n)
			}
			return nil
		}),
		customdiff.ValidateChange("fhir_store", func(ctx context.Context, old, new, m interface{}) error {
			o := old.(string)
			n := new.(string)
			if o == "" { // New
				return nil
			}
			if o != n {
				return fmt.Errorf("changing the fhir_store endpoint from '%s' to '%s' is not supported", o, n)
			}
			return nil
		}),
		customdiff.ValidateChange("org_id", func(ctx context.Context, old, new, m interface{}) error {
			o := old.(string)
			n := new.(string)
			if o == "" { // New
				return nil
			}
			if o != n {
				return fmt.Errorf("changing the organization id from '%s' to '%s' is not supported", o, n)
			}
			return nil
		}),
	)
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

	if d.HasChange("version") || d.HasChange("fhir_store") || d.HasChange("org_id") {
		return diag.FromErr(fmt.Errorf("changes in 'version', 'fhir_stor' or 'org_id' are not supported"))
	}

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
