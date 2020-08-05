package hsdp

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
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

	var client *creds.Client
	var err error

	if username != "" && password != "" {
		client, err = config.CredentialsClientWithLogin(username, password)
		if err != nil {
			return err
		}
	} else {
		client, err = config.CredentialsClient()
		if err != nil {
			return err
		}
	}
	if client == nil {
		return ErrMissingClientPassword
	}
	credentials, _, err := client.Access.GetAccess(&creds.GetAccessOptions{
		ProductKey: &productKey,
	})
	if err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(&credentials)
	if err != nil {
		return err
	}
	d.SetId("access")
	d.Set("access", string(jsonBytes))

	return err
}
