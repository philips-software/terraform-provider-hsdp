package practitioner

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func DataSourceCDRPractitioner() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCDRPropositionRead,
		Schema: map[string]*schema.Schema{
			"fhir_store": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fhir_json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity_systems": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     tools.StringSchema(),
			},
			"identity_values": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     tools.StringSchema(),
			},
			"identity_uses": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     tools.StringSchema(),
			},
		},
	}

}

func dataSourceCDRPropositionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)

	fhirStore := d.Get("fhir_store").(string)

	client, err := c.GetFHIRClientFromEndpoint(fhirStore)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	version := d.Get("version").(string)
	guid := d.Get("guid").(string)

	d.SetId(guid)
	switch version {
	case "r4":
		contained, resp, err := client.OperationsR4.Get("Practitioner/" + d.Id())
		if err != nil {
			if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone) {
				d.SetId("")
				return diags
			}
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Errorf("practitioner read: %w", err).Error(),
				})
			}
			return diags
		}
		resource := contained.GetPractitioner()
		jsonResource, err := c.R4MA.MarshalResource(resource)
		if err != nil {
			return diag.FromErr(fmt.Errorf("R4.MarshalResource: %w", err))
		}
		_ = d.Set("fhir_json", string(jsonResource))
	case "stu3", "":
		contained, resp, err := client.OperationsSTU3.Get("Practitioner/" + d.Id())
		if err != nil {
			if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone) {
				d.SetId("")
				return diags
			}
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Errorf("practitioner read: %w", err).Error(),
				})
			}
			return diags
		}
		resource := contained.GetPractitioner()
		jsonResource, err := c.STU3MA.MarshalResource(resource)
		if err != nil {
			return diag.FromErr(fmt.Errorf("STU3.MarshalResource: %w", err))
		}
		_ = d.Set("fhir_json", string(jsonResource))
	default:
		return diag.FromErr(fmt.Errorf("unsupported FHIR version '%s'", version))
	}
	return diags
}
