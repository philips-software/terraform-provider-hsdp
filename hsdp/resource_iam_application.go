package hsdp

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
	"net/http"
)

func resourceIAMApplication() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceIAMApplicationCreate,
		ReadContext:   resourceIAMApplicationRead,
		DeleteContext: resourceIAMApplicationDelete,

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
			"proposition_id": {
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

func resourceIAMApplicationCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var app iam.Application
	app.Name = d.Get("name").(string) // TODO: this must be all caps
	app.Description = d.Get("description").(string)
	app.PropositionID = d.Get("proposition_id").(string)
	app.GlobalReferenceID = d.Get("global_reference_id").(string)

	createdApp, resp, err := client.Applications.CreateApplication(app)
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode != http.StatusConflict {
			return diag.FromErr(err)
		}
		createdApp, resp, err = client.Applications.GetApplicationByName(app.Name)
		if err != nil {
			return diag.FromErr(err)
		}
		if createdApp.Description != app.Description {
			return diag.FromErr(fmt.Errorf("existing application found but description mismatch: '%s' != '%s'", createdApp.Description, app.Description))
		}
		if createdApp.PropositionID != app.PropositionID {
			return diag.FromErr(fmt.Errorf("existing application found but proposition_id mismatch: '%s' != '%s'", createdApp.PropositionID, app.PropositionID))
		}
		if createdApp.GlobalReferenceID != app.GlobalReferenceID {
			return diag.FromErr(fmt.Errorf("existing application found but global_reference_id mismatch: '%s' != '%s'", createdApp.GlobalReferenceID, app.GlobalReferenceID))
		}
		// We found a matching existing application, go with it
	}
	if createdApp == nil {
		return diag.FromErr(fmt.Errorf("Unexpected failure creating '%s': [%v] [%v]", app.Name, err, resp))
	}
	d.SetId(createdApp.ID)
	_ = d.Set("name", createdApp.Name)
	_ = d.Set("description", createdApp.Description)
	_ = d.Set("proposition_id", createdApp.PropositionID)
	_ = d.Set("global_reference_id", createdApp.GlobalReferenceID)
	return diags
}

func resourceIAMApplicationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	app, resp, err := client.Applications.GetApplicationByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("name", app.Name)
	_ = d.Set("description", app.Description)
	_ = d.Set("proposition_id", app.PropositionID)
	_ = d.Set("global_reference_id", app.GlobalReferenceID)
	return diags
}

func resourceIAMApplicationUpdate(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if !d.HasChange("description") {
		return diags
	}
	// Not implemented by HSDP
	return diags
}

func resourceIAMApplicationDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}
