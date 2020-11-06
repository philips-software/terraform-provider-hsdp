package hsdp

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	creds "github.com/philips-software/go-hsdp-api/credentials"
)

func dataSourceS3CredentialsAccess() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceS3CredsAccessRead,
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

func dataSourceS3CredsAccessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var diags diag.Diagnostics

	productKey := d.Get("product_key").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	var client *creds.Client
	var err error

	if username != "" && password != "" {
		client, err = config.CredentialsClientWithLogin(username, password)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		client, err = config.CredentialsClient()
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if client == nil {
		return diag.FromErr(ErrMissingClientPassword)
	}
	credentials, _, err := client.Access.GetAccess(&creds.GetAccessOptions{
		ProductKey: &productKey,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	jsonBytes, err := json.Marshal(&credentials)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("access")
	d.Set("access", string(jsonBytes))

	return diags

}
