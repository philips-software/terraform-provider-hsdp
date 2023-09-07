package blr

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceBLRBlobStorePolicyDefinition() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBLRStorePolicyRead,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"principal": config.PrincipalSchema(),
		},
	}

}

func dataSourceBLRStorePolicyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	principal := config.SchemaToPrincipal(d, m)

	policyID := d.Get("policy_id").(string)

	client, err := c.BLRClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	resource, _, err := client.Configurations.GetBlobStorePolicyByID(policyID)
	if err != nil {
		return diag.FromErr(err)
	}

	b, err := json.Marshal(resource.Statement)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(policyID)
	_ = d.Set("policy", string(b))

	return diags
}
