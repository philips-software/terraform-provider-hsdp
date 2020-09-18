package hsdp

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func dataSourceIAMApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIAMApplicationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"proposition_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"global_reference_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceIAMApplicationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}
	propID := d.Get("proposition_id").(string)
	name := d.Get("name").(string)

	prop, _, err := client.Applications.GetApplication(&iam.GetApplicationsOptions{
		PropositionID: &propID,
		Name:          &name,
	})
	if err != nil {
		return err
	}
	if len(prop) == 0 {
		return ErrResourceNotFound
	}

	d.SetId(prop[0].ID)
	_ = d.Set("description", prop[0].Description)
	_ = d.Set("global_reference_id", prop[0].GlobalReferenceID)
	return nil
}
