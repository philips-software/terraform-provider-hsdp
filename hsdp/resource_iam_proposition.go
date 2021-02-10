package hsdp

import (
	"context"
	"fmt"
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
		DeleteContext: resourceIAMPropositionDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateUpperString,
				ForceNew:     true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"global_reference_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIAMPropositionCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	createdProp, resp, err := client.Propositions.CreateProposition(prop)
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode != http.StatusConflict {
			return diag.FromErr(err)
		}
		createdProp, _, err = client.Propositions.GetProposition(&iam.GetPropositionsOptions{
			Name: &prop.Name,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		if createdProp.Description != prop.Description {
			return diag.FromErr(fmt.Errorf("existing proposition found but description mismatch: '%s' != '%s'", createdProp.Description, prop.Description))
		}
		if createdProp.OrganizationID != prop.OrganizationID {
			return diag.FromErr(fmt.Errorf("existing proposition found but organization_id mismatch: '%s' != '%s'", createdProp.OrganizationID, prop.OrganizationID))
		}
		if createdProp.GlobalReferenceID != prop.GlobalReferenceID {
			return diag.FromErr(fmt.Errorf("existing proposition found but global_reference_id mismatch: '%s' != '%s'", createdProp.OrganizationID, prop.OrganizationID))
		}
		// We found a matching existing proposition, go with it
	}
	d.SetId(createdProp.ID)
	_ = d.Set("name", createdProp.Name)
	_ = d.Set("description", createdProp.Description)
	_ = d.Set("organization_id", createdProp.OrganizationID)
	_ = d.Set("global_reference_id", createdProp.GlobalReferenceID)
	return diags
}

func resourceIAMPropositionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceIAMPropositionUpdate(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if !d.HasChange("description") {
		return diags
	}
	return diag.FromErr(ErrNotImplementedByHSDP)
}

func resourceIAMPropositionDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// As HSDP IAM does not support IAM proposition deletion we simply
	// clear the proposition from state. This will be properly implemented
	// once the IAM API balances out
	var diags diag.Diagnostics
	d.SetId("")

	return diags
}
