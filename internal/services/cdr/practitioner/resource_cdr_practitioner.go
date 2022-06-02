package practitioner

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceCDRPractitioner() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: importPractitionerContext,
		},

		CreateContext: resourceCDRPractitionerCreate,
		ReadContext:   resourceCDRPractitionerRead,
		UpdateContext: resourceCDRPractitionerUpdate,
		DeleteContext: resourceCDRPractitionerDelete,

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
			"identifier": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     identifierSchema(),
			},
			"name": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     nameSchema(),
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func nameSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"text": {
				Type:     schema.TypeString,
				Required: true,
			},
			"family": {
				Type:     schema.TypeString,
				Required: true,
			},
			"given": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     tools.StringSchema(),
			},
		},
	}
}

func identifierSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"system": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCDRPractitionerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	fhirStore := d.Get("fhir_store").(string)

	client, err := c.GetFHIRClientFromEndpoint(fhirStore)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	version := d.Get("version").(string)

	switch version {
	case "stu3":
		if createDiags := stu3Create(ctx, c, client, d, m); len(createDiags) > 0 {
			return createDiags
		}
	case "r4":
		if createDiags := r4Create(ctx, c, client, d, m); len(createDiags) > 0 {
			return createDiags
		}
	default:
		return diag.FromErr(fmt.Errorf("unsupported FHIR version '%s'", version))
	}
	return resourceCDRPractitionerRead(ctx, d, m)
}

func resourceCDRPractitionerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	fhirStore := d.Get("fhir_store").(string)
	version := d.Get("version").(string)

	if fhirStore == "" {
		return diag.FromErr(fmt.Errorf("practitioner read: the 'fhir_store' attribute is blank"))
	}

	client, err := c.GetFHIRClientFromEndpoint(fhirStore)
	if err != nil {
		return diag.FromErr(fmt.Errorf("practitioner read: %w", err))
	}
	defer client.Close()

	switch version {
	case "stu3":
		if readDiags := stu3Read(ctx, c, client, d, m); len(readDiags) > 0 {
			return readDiags
		}
	case "r4":
		if readDiags := r4Read(ctx, c, client, d, m); len(readDiags) > 0 {
			return readDiags
		}
	default:
		return diag.FromErr(fmt.Errorf("unsupported FHIR version '%s'", version))
	}
	return diags
}

func resourceCDRPractitionerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics

	fhirStore := d.Get("fhir_store").(string)
	version := d.Get("version").(string)

	client, err := c.GetFHIRClientFromEndpoint(fhirStore)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	switch version {
	case "stu3":
		if readDiags := stu3Update(ctx, c, client, d, m); len(readDiags) > 0 {
			return readDiags
		}
	case "r4":
		if readDiags := r4Update(ctx, c, client, d, m); len(readDiags) > 0 {
			return readDiags
		}
	default:
		return diag.FromErr(fmt.Errorf("unsupported FHIR version '%s'", version))
	}
	return diags
}

func resourceCDRPractitionerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics

	fhirStore := d.Get("fhir_store").(string)
	version := d.Get("version").(string)

	client, err := c.GetFHIRClientFromEndpoint(fhirStore)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	switch version {
	case "stu3":
		if readDiags := stu3Delete(ctx, c, client, d, m); len(readDiags) > 0 {
			return readDiags
		}
	case "r4":
		if readDiags := r4Delete(ctx, c, client, d, m); len(readDiags) > 0 {
			return readDiags
		}
	default:
		return diag.FromErr(fmt.Errorf("unsupported FHIR version '%s'", version))
	}
	return diags

}
