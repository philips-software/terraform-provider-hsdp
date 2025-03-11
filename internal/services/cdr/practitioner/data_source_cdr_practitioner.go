package practitioner

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dip-software/go-dip-api/cdr"
	r4idhelper "github.com/dip-software/go-dip-api/cdr/helper/fhir/r4/identifier"
	stu3idhelper "github.com/dip-software/go-dip-api/cdr/helper/fhir/stu3/identifier"
	r4pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/stu3/resources_go_proto"
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
				Type:     schema.TypeList,
				Computed: true,
				Elem:     tools.StringSchema(),
			},
			"identity_values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     tools.StringSchema(),
			},
			"identity_uses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     tools.StringSchema(),
			},
		},
	}

}

func dataSourceCDRPropositionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		var contained *r4pb.ContainedResource
		var resp *cdr.Response
		err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
			var err error
			contained, resp, err = client.OperationsR4.Get("Practitioner/" + d.Id())
			if err != nil {
				_ = client.TokenRefresh()
			}
			if resp == nil {
				return nil, err
			}
			return resp.Response, err
		})
		if err != nil {
			if resp != nil && (resp.StatusCode() == http.StatusNotFound || resp.StatusCode() == http.StatusGone) {
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
		if len(resource.Identifier) > 0 {
			var uses []string
			var systems []string
			var values []string
			for _, i := range resource.Identifier {
				uses = append(uses, r4idhelper.UseToString(i.Use))
				systems = append(systems, i.System.GetValue())
				values = append(values, i.Value.GetValue())
			}
			_ = d.Set("identity_uses", uses)
			_ = d.Set("identity_systems", systems)
			_ = d.Set("identity_values", values)
		}
	case "stu3", "":
		var contained *resources_go_proto.ContainedResource
		var resp *cdr.Response
		err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
			var err error
			contained, resp, err = client.OperationsSTU3.Get("Practitioner/" + d.Id())
			if err != nil {
				_ = client.TokenRefresh()
			}
			if resp == nil {
				return nil, err
			}
			return resp.Response, err
		})
		if err != nil {
			if resp != nil && (resp.StatusCode() == http.StatusNotFound || resp.StatusCode() == http.StatusGone) {
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
		if len(resource.Identifier) > 0 {
			var uses []string
			var systems []string
			var values []string
			for _, i := range resource.Identifier {
				uses = append(uses, stu3idhelper.UseToString(i.Use))
				systems = append(systems, i.System.String())
				values = append(values, i.Value.String())
			}
			_ = d.Set("identity_uses", uses)
			_ = d.Set("identity_systems", systems)
			_ = d.Set("identity_values", values)
		}
	default:
		return diag.FromErr(fmt.Errorf("unsupported FHIR version '%s'", version))
	}
	return diags
}
