package hsdp

import (
	"errors"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func resourceIAMProposition() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceIAMPropositionCreate,
		Read:   resourceIAMPropositionRead,
		Update: resourceIAMPropositionUpdate,
		Delete: resourceIAMPropositionDelete,

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

func resourceIAMPropositionCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	var prop iam.Proposition
	prop.Name = d.Get("name").(string) // TODO: this must be all caps
	prop.Description = d.Get("description").(string)
	prop.OrganizationID = d.Get("organization_id").(string)
	prop.GlobalReferenceID = d.Get("global_reference_id").(string)

	createdProp, _, err := client.Propositions.CreateProposition(prop)
	if err != nil {
		return err
	}
	d.SetId(createdProp.ID)
	d.Set("name", createdProp.Name)
	d.Set("description", createdProp.Description)
	d.Set("organization_id", createdProp.OrganizationID)
	d.Set("global_reference_id", createdProp.GlobalReferenceID)
	return nil
}

func resourceIAMPropositionRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	id := d.Id()
	prop, resp, err := client.Propositions.GetPropositionByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("name", prop.Name)
	d.Set("description", prop.Description)
	d.Set("organization_id", prop.OrganizationID)
	d.Set("global_reference_id", prop.GlobalReferenceID)
	return nil
}

func resourceIAMPropositionUpdate(d *schema.ResourceData, m interface{}) error {
	if !d.HasChange("description") {
		return nil
	}
	return errors.New("not implemented by HSDP")
}

func resourceIAMPropositionDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
