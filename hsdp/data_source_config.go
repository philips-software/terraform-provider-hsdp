package hsdp

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/philips-software/go-hsdp-api/config"
)

func dataSourceConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConfigRead,
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

func dataSourceConfigRead(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(*Config)

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
		return err
	}
	d.SetId("data" + region + environment + service)
	if url, err := c.Service(service).GetString("url"); err == nil {
		_ = d.Set("url", url)
	}
	if host, err := c.Service(service).GetString("host"); err == nil {
		_ = d.Set("host", host)
	}
	if domain, err := c.Service(service).GetString("domain"); err == nil {
		_ = d.Set("domain", domain)
	}
	return nil
}
