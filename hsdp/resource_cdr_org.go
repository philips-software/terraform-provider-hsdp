package hsdp

import (
	"context"
	"github.com/google/fhir/go/proto/google/fhir/proto/stu3/datatypes_go_proto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	jsonpatch "github.com/herkyl/patchwerk"
	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/stu3"
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
		},
	}
}

func resourceCDROrgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	endpoint := d.Get("fhir_store").(string)
	orgID := d.Get("org_id").(string)

	client, err := config.getFHIRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	org, err := stu3.NewOrganization(config.TimeZone, orgID, name)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	partOf := d.Get("part_of").(string)
	if partOf != "" {
		org.PartOf = &datatypes_go_proto.Reference{
			Reference: &datatypes_go_proto.Reference_OrganizationId{
				OrganizationId: &datatypes_go_proto.ReferenceId{
					Value: partOf,
				},
			},
		}
	}

	// Check if already onboarded
	onboardedOrg, _, err := client.TenantSTU3.GetOrganizationByID(orgID)
	if err == nil && onboardedOrg != nil {
		d.SetId(onboardedOrg.Id.Value)
		return resourceCDROrgUpdate(ctx, d, m)
	}
	// Do initial boarding
	onboardedOrg, _, err = client.TenantSTU3.Onboard(org)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(onboardedOrg.Id.Value)
	return diags
}

func resourceCDROrgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	endpoint := d.Get("fhir_store").(string)
	orgID := d.Get("org_id").(string)

	client, err := config.getFHIRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
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
	if org.PartOf != nil {
		_ = d.Set("part_of", org.PartOf.GetOrganizationId())
	}
	return diags
}

func resourceCDROrgUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	var diags diag.Diagnostics

	endpoint := d.Get("fhir_store").(string)
	id := d.Id()

	client, err := config.getFHIRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	org, _, err := client.TenantSTU3.GetOrganizationByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonOrg, err := config.ma.MarshalResource(org)
	if err != nil {
		return diag.FromErr(err)
	}
	madeChanges := false

	if d.HasChange("name") {
		org.Name.Value = d.Get("name").(string)
		madeChanges = true
	}
	if d.HasChange("part_of") {
		partOf := d.Get("part_of").(string)
		if partOf != "" {
			org.PartOf = &datatypes_go_proto.Reference{
				Reference: &datatypes_go_proto.Reference_OrganizationId{
					OrganizationId: &datatypes_go_proto.ReferenceId{
						Value: partOf,
					},
				},
			}
		} else {
			org.PartOf = nil
		}
		madeChanges = true
	}
	if !madeChanges {
		return diags
	}

	changedOrg, _ := config.ma.MarshalResource(org)
	patch, err := jsonpatch.DiffBytes(jsonOrg, changedOrg)
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = client.OperationsSTU3.Patch("Organization/"+id, patch)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCDROrgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// TODO: will be supported in CDR release of Q1 2021
	return diags
}
