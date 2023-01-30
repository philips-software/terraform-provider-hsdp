package subscription

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceCDRSubscription() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: importSubscriptionContext,
		},

		Description:   "Provides a resource for managing FHIR Subscriptions in a CDR. The only supported channel type is rest-webhook therefore the endpoint and headers are top-level arguments.",
		CreateContext: resourceCDRSubscriptionCreate,
		ReadContext:   resourceCDRSubscriptionRead,
		UpdateContext: resourceCDRSubscriptionUpdate,
		DeleteContext: resourceCDRSubscriptionDelete,
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceCDRSubscriptionV0().CoreConfigSchema().ImpliedType(),
				Upgrade: patchSubscriptionV0,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"fhir_store": {
				Description: "The CDR FHIR store tenant URL.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"version": {
				Description: "The FHIR version to use. Options: [ 'stu3' | 'r4' ].",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "stu3",
				ForceNew:    true,
			},
			"delete_endpoint": {
				Description: "The REST endpoint to call for DELETE operations. Must use https schema.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"criteria": {
				Description: "On which resource to notify.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"endpoint": {
				Description: "The REST endpoint to call. Must use https schema.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"reason": {
				Description: "Reason for creating the subscription.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"headers": {
				Description: "List of headers to add to the REST call.",
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"end": {
				Description: "RFC3339 formatted timestamp when to end notifications.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"status": {
				Description: "The status of the subscription (requested | active | error | off).",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceCDRSubscriptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	return resourceCDRSubscriptionRead(ctx, d, m)
}

func resourceCDRSubscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	fhirStore := d.Get("fhir_store").(string)
	version := d.Get("version").(string)

	if fhirStore == "" {
		return diag.FromErr(fmt.Errorf("subscription read: the 'fhir_store' attribute is blank"))
	}

	client, err := c.GetFHIRClientFromEndpoint(fhirStore)
	if err != nil {
		return diag.FromErr(fmt.Errorf("subscription read: %w", err))
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

func resourceCDRSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceCDRSubscriptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
