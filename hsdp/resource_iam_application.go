package hsdp

import (
	"context"
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
		UpdateContext: resourceIAMApplicationUpdate,
		DeleteContext: resourceIAMApplicationDelete,

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
			"proposition_id": &schema.Schema{
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

func resourceIAMApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	createdApp, _, err := client.Applications.CreateApplication(app)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(createdApp.ID)
	_ = d.Set("name", createdApp.Name)
	_ = d.Set("description", createdApp.Description)
	_ = d.Set("proposition_id", createdApp.PropositionID)
	_ = d.Set("global_reference_id", createdApp.GlobalReferenceID)
	return diags
}

func resourceIAMApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceIAMApplicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if !d.HasChange("description") {
		return diags
	}
	// Not implemented by HSDP
	return diags
}

func resourceIAMApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
