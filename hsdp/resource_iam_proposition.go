package hsdp

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func resourceIAMProposition() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceIAMPropositionCreate,
		ReadContext:   resourceIAMPropositionRead,
		UpdateContext: resourceIAMPropositionUpdate,
		DeleteContext: resourceIAMPropositionDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateUpperString,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"global_reference_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceIAMPropositionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var prop iam.Proposition
	prop.Name = d.Get("name").(string) // TODO: this must be all caps
	prop.Description = d.Get("description").(string)
	prop.OrganizationID = d.Get("organization_id").(string)
	prop.GlobalReferenceID = d.Get("global_reference_id").(string)

	createdProp, _, err := client.Propositions.CreateProposition(prop)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(createdProp.ID)
	_ = d.Set("name", createdProp.Name)
	_ = d.Set("description", createdProp.Description)
	_ = d.Set("organization_id", createdProp.OrganizationID)
	_ = d.Set("global_reference_id", createdProp.GlobalReferenceID)
	return diags
}

func resourceIAMPropositionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	prop, resp, err := client.Propositions.GetPropositionByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("name", prop.Name)
	_ = d.Set("description", prop.Description)
	_ = d.Set("organization_id", prop.OrganizationID)
	_ = d.Set("global_reference_id", prop.GlobalReferenceID)
	return diags
}

func resourceIAMPropositionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if !d.HasChange("description") {
		return diags
	}
	return diag.FromErr(ErrNotImplementedByHSDP)
}

func resourceIAMPropositionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// As HSDP IAM does not support IAM proposition deletion we simply
	// clear the proposition from state. This will be properly implemented
	// once the IAM API balances out
	var diags diag.Diagnostics

	return diags
}
