package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func dataSourceIAMIntrospect() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIAMIntrospectRead,
		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"managing_organization": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceIAMIntrospectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*iam.Client)

	resp, _, err := client.Introspect()

	if err != nil {
		return err
	}
	d.Set("managing_organization", resp.Organizations.ManagingOrganization)
	d.SetId(resp.Username)
	d.Set("username", resp.Username)
	d.Set("token", client.Token())

	return err
}
