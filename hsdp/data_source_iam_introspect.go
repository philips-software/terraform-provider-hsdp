package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hsdp/go-hsdp-iam/api"
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
	client := meta.(*api.Client)

	resp, _, err := client.Introspect()

	d.Set("managing_organization", resp.Organizations.ManagingOrganization)
	d.SetId(resp.Username)
	d.Set("username", resp.Username)
	d.Set("token", client.Token())

	return err
}
