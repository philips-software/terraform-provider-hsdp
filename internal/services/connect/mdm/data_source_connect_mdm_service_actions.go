package mdm

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceConnectMDMServiceActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMServiceActionsRead,
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
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"guids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"organization_guid_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"standard_service_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}

}

func dataSourceConnectMDMServiceActionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	name := ""
	orgGuidValue := ""
	standardServiceId := ""
	id := ""

	if v, ok := d.GetOk("filter"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			name = mVi["name"].(string)
			orgGuidValue = mVi["organization_guid_value"].(string)
			standardServiceId = mVi["standard_service_id"].(string)
			id = mVi["id"].(string)
		}
	}
	findOpts := &mdm.GetServiceActionOptions{}
	if name != "" {
		findOpts.Name = &name
	}
	if orgGuidValue != "" {
		findOpts.OrganizationGuidValue = &orgGuidValue
	}
	if standardServiceId != "" {
		findOpts.StandardServiceID = &standardServiceId
	}
	if id != "" {
		findOpts.ID = &id
	}

	resources, _, err := client.ServiceActions.Find(findOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []string
	var names []string
	var descriptions []string
	var guids []string
	var types []string

	for _, r := range *resources {
		ids = append(ids, r.ID)
		names = append(names, r.Name)
		descriptions = append(descriptions, r.Description)
		guids = append(guids, guidValue(r))
		types = append(types, r.ResourceType)
	}
	_ = d.Set("ids", ids)
	_ = d.Set("names", names)
	_ = d.Set("descriptions", descriptions)
	_ = d.Set("guids", guids)
	_ = d.Set("types", types)

	result, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(result)
	return diags
}
