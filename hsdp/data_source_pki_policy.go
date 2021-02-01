package hsdp

import (
	"bytes"
	"context"
	"encoding/pem"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/pki"
)

func dataSourcePKIPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePKIPolicyRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ca_pem": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"crl_pem": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourcePKIPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	var diags diag.Diagnostics
	var err error
	var client *pki.Client

	region := d.Get("region").(string)
	environment := d.Get("environment").(string)
	if region != "" || environment != "" {
		client, err = config.PKIClient(region, environment)
	} else {
		client, err = config.PKIClient()
	}
	if err != nil {
		return diag.FromErr(err)
	}

	// Policy CA
	ca, block, _, err := client.Services.GetPolicyCA()
	if err != nil {
		return diag.FromErr(err)
	}
	var caPem bytes.Buffer
	if err := pem.Encode(&caPem, block); err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("ca_pem", caPem.String())

	// Policy CRL
	_, block, _, err = client.Services.GetPolicyCRL()
	if err != nil {
		return diag.FromErr(err)
	}
	var crlPem bytes.Buffer
	if err := pem.Encode(&crlPem, block); err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("crl_pem", crlPem.String())

	d.SetId(fmt.Sprintf("%v", ca.SerialNumber))
	return diags
}
