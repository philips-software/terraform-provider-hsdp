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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region to look up. Defaults to the provider configured one",
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The environment to refer to. Defaults to the provider configured one",
			},
			"service": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service to look up",
			},
			"host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host associated with the service",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The url associated with the service",
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain associated with the service",
			},
			"services": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The service ID used to authenticate against IAM",
			},
			"org_admin_username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OrgAdmin username used to authenticate against IAM",
			},
		},
	}

}

func dataSourceConfigRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerConfig := m.(*Config)

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
	_ = d.Set("services", c.Services())
	_ = d.Set("service_id", providerConfig.ServiceID)
	_ = d.Set("org_admin_username", providerConfig.OrgAdminUsername)

	return diags
}
