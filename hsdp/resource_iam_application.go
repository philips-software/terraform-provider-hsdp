package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
	"net/http"
)

func resourceIAMApplication() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceIAMApplicationCreate,
		Read:   resourceIAMApplicationRead,
		Update: resourceIAMApplicationUpdate,
		Delete: resourceIAMApplicationDelete,

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

func resourceIAMApplicationCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	var app iam.Application
	app.Name = d.Get("name").(string) // TODO: this must be all caps
	app.Description = d.Get("description").(string)
	app.PropositionID = d.Get("proposition_id").(string)
	app.GlobalReferenceID = d.Get("global_reference_id").(string)

	createdApp, _, err := client.Applications.CreateApplication(app)
	if err != nil {
		return err
	}
	d.SetId(createdApp.ID)
	d.Set("name", createdApp.Name)
	d.Set("description", createdApp.Description)
	d.Set("proposition_id", createdApp.PropositionID)
	d.Set("global_reference_id", createdApp.GlobalReferenceID)
	return nil
}

func resourceIAMApplicationRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	id := d.Id()
	app, resp, err := client.Applications.GetApplicationByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("name", app.Name)
	d.Set("description", app.Description)
	d.Set("proposition_id", app.PropositionID)
	d.Set("global_reference_id", app.GlobalReferenceID)
	return nil
}

func resourceIAMApplicationUpdate(d *schema.ResourceData, m interface{}) error {
	if !d.HasChange("description") {
		return nil
	}
	// Not implemented by HSDP
	return nil
}

func resourceIAMApplicationDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
