package hsdp

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func dataSourceIAMProposition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIAMPropositionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
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

func dataSourceIAMPropositionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}
	orgId := d.Get("organization_id").(string)
	name := d.Get("name").(string)

	prop, _, err := client.Propositions.GetProposition(&iam.GetPropositionsOptions{
		OrganizationID: &orgId,
		Name:           &name,
	})
	if err != nil {
		return err
	}

	d.SetId(orgId)
	_ = d.Set("description", prop.Description)
	_ = d.Set("global_reference_id", prop.GlobalReferenceID)
	return nil
}
