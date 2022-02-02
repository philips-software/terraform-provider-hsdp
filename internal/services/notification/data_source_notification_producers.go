package notification

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/notification"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceNotificationProducers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationProducersRead,
		Schema: map[string]*schema.Schema{
			"managing_organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"producer_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceNotificationProducersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	managingOrgID := d.Get("managing_organization_id").(string)

	opts := &notification.GetOptions{
		ManagedOrganizationID: &managingOrgID,
	}

	list, resp, err := client.Producer.GetProducers(opts) // Get all producers

	if err != nil {
		if resp == nil || resp.StatusCode != http.StatusForbidden { // Do not error on permission issues
			return diag.FromErr(err)
		}
		list = []notification.Producer{} // empty list
	}
	producers := make([]string, 0)

	for _, p := range list {
		producers = append(producers, p.ID)
	}
	d.SetId(fmt.Sprintf("notification-producers-%s", strings.Join(producers, "-")))
	_ = d.Set("producer_ids", producers)

	return diags
}
