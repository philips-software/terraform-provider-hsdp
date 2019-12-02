package hsdp

import (
	"encoding/json"

	"github.com/hashicorp/terraform/helper/schema"
	creds "github.com/philips-software/go-hsdp-api/credentials"
)

func dataSourceS3CredentialsAccess() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceS3CredsAccessRead,
		Schema: map[string]*schema.Schema{
			"access": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}

}

func dataSourceS3CredsAccessRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	productKey := d.Get("product_key").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	client, err := config.CredentialsClientWithLogin(username, password)
	if err != nil {
		return err
	}

	creds, _, err := client.Access.GetAccess(&creds.GetAccessOptions{
		ProductKey: &productKey,
	})
	if err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(&creds)
	if err != nil {
		return err
	}
	d.SetId("access")
	d.Set("access", string(jsonBytes))

	return err
}
