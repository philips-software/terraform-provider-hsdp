package hsdp

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/config"
)

func dataSourceConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*Config)

	var diags diag.Diagnostics

	service := d.Get("service").(string)
	region := d.Get("region").(string)
	environment := d.Get("environment").(string)
	if region == "" {
		region = providerConfig.Region
	}
	if environment == "" {
		environment = providerConfig.Environment
	}
	c, err := config.New(config.WithEnv(environment),
		config.WithRegion(region))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("data" + region + environment + service)
	if url := c.Service(service).URL; url != "" {
		_ = d.Set("url", url)
	}
	if host := c.Service(service).Host; host != "" {
		_ = d.Set("host", host)
	}
	if domain := c.Service(service).Domain; domain != "" {
		_ = d.Set("domain", domain)
	}
	return diags
}
