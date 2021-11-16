package mdm

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceConnectMDMStandardServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMStandardServicesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"trusted": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
			},
		},
	}

}

func dataSourceConnectMDMStandardServicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	standardServices, _, err := client.StandardServices.GetStandardServices(nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []string
	var names []string
	var descriptions []string
	var trusted []bool

	for _, s := range *standardServices {
		// All criteria match, so add user
		ids = append(ids, s.ID)
		names = append(names, s.Name)
		descriptions = append(descriptions, s.Description)
		trusted = append(trusted, s.Trusted)
	}
	_ = d.Set("ids", ids)
	_ = d.Set("names", names)
	_ = d.Set("descriptions", descriptions)
	_ = d.Set("trusted", trusted)
	result, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(result)
	return diags
}
