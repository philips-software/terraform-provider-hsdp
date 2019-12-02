package hsdp

import (
	"encoding/json"

	"github.com/hashicorp/terraform/helper/schema"
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
			"introspect": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceIAMIntrospectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.IAMClient()

	resp, _, err := client.Introspect()

	if err != nil {
		return err
	}
	introspectJSON, err := json.Marshal(&resp)
	if err != nil {
		return err
	}

	d.Set("managing_organization", resp.Organizations.ManagingOrganization)
	d.SetId(resp.Username)
	d.Set("username", resp.Username)
	d.Set("token", client.Token())
	d.Set("introspect", string(introspectJSON))

	return err
}
