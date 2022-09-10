package discovery

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/discovery"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func DataSourceDiscoveryService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDiscoveryServiceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"is_trusted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"urls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"principal": config.PrincipalSchema(),
		},
	}

}

func dataSourceDiscoveryServiceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.DiscoveryClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	tag := d.Get("tag").(string)

	var services *[]discovery.Service

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var resp *discovery.Response
		services, resp, err = client.GetServices()
		if resp == nil {
			return nil, fmt.Errorf("response is nil: %w", err)
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("GetServices: %w", err))
	}
	var service *discovery.Service

	if name != "" {
		for _, s := range *services {
			if s.Name == name {
				service = &s
				break
			}
		}
	}
	if tag != "" {
		for _, s := range *services {
			if s.Tag == tag {
				service = &s
				break
			}
		}
	}
	if service == nil {
		return diag.FromErr(fmt.Errorf("no service found matching name '%s' or tag '%s'", name, tag))
	}
	_ = d.Set("actions", service.Actions)
	_ = d.Set("urls", service.URLS)
	_ = d.Set("name", service.Name)
	_ = d.Set("tag", service.Tag)
	_ = d.Set("is_trusted", service.IsTrusted)

	d.SetId(service.ID)
	return diags
}
