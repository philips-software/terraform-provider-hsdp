package mdm

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceConnectMDMDataSubscribers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMDataSubscribersRead,
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
			"configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subscriber_guids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subscriber_type_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceConnectMDMDataSubscribersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resources, _, err := client.DataSubscribers.Get(nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []string
	var names []string
	var configurations []string
	var subscriberGuids []string
	var subscriberTypeIds []string

	for _, r := range *resources {
		ids = append(ids, r.ID)
		names = append(names, r.Name)
		configurations = append(configurations, string(r.Configuration))
		subscriberGuids = append(subscriberGuids, fmt.Sprintf("%s|%s", r.SubscriberGuid.System, r.SubscriberGuid.Value))
		subscriberTypeIds = append(subscriberTypeIds, r.SubscriberTypeId.Reference)
	}
	_ = d.Set("ids", ids)
	_ = d.Set("names", names)
	_ = d.Set("configurations", configurations)
	_ = d.Set("subscriber_guids", subscriberGuids)
	_ = d.Set("subscriber_type_ids", subscriberTypeIds)

	result, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(result)
	return diags
}
