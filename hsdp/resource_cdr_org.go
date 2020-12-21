package hsdp

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/stu3"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCDROrg() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceCDROrgCreate,
		ReadContext:   resourceCDROrgRead,
		UpdateContext: resourceCDROrgUpdate,
		DeleteContext: resourceCDROrgDelete,

		Schema: map[string]*schema.Schema{
			"fhir_store": {
				Type:     schema.TypeString,
				Required: true,
			},
			"root_org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCDROrgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	fhirStore := d.Get("fhir_store").(string)
	rootOrgID := d.Get("root_org_id").(string)
	orgID := d.Get("org_id").(string)

	client, err := config.getFHIRClient(fhirStore, rootOrgID)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	org, err := stu3.NewOrganization("Europe/Amsterdam", orgID, name)
	if err != nil {
		return diag.FromErr(err)
	}

	onboardedOrg, resp, err := client.TenantSTU3.Onboard(org)
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode != http.StatusConflict {
			return diag.FromErr(err)
		}
		onboardedOrg, _, err = client.TenantSTU3.GetOrganizationByID(orgID)
		if err != nil {
			return diag.FromErr(err)
		}
		// We found a matching existing org, go with it
	}
	d.SetId(onboardedOrg.Id.Value)
	return diags
}

func resourceCDROrgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	fhirStore := d.Get("fhir_store").(string)
	rootOrgID := d.Get("root_org_id").(string)
	orgID := d.Get("org_id").(string)

	client, err := config.getFHIRClient(fhirStore, rootOrgID)
	if err != nil {
		return diag.FromErr(err)
	}
	org, resp, err := client.TenantSTU3.GetOrganizationByID(orgID)
	if err != nil || resp == nil {
		if resp == nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "response is nil",
			})
		}
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})
		}
		return diags
	}
	_ = d.Set("name", org.Name.Value)
	return diags
}

func resourceCDROrgUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// TODO: no meaningful lifecycle here
	if !d.HasChange("name") {
		return diags
	}
	return diag.FromErr(ErrNotImplementedByHSDP)
}

func resourceCDROrgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// TODO
	return diags
}
