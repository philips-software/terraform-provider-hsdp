package hsdp

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	resp, _, err := client.Introspect()

	if err != nil {
		return err
	}
	introspectJSON, err := json.Marshal(&resp)
	if err != nil {
		return err
	}

	d.SetId(resp.Username)
	_ = d.Set("managing_organization", resp.Organizations.ManagingOrganization)
	_ = d.Set("username", resp.Username)
	_ = d.Set("token", client.Token())
	_ = d.Set("introspect", string(introspectJSON))

	return err
}
