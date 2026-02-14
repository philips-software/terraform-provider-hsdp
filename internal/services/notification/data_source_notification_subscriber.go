package notification

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-dip-api/notification"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceNotificationSubscriber() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationSubscriberRead,
		Schema: map[string]*schema.Schema{
			"subscriber_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"managing_organization": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscriber_product_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subscriber_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscriber_service_instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscriber_service_base_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscriber_service_path_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNotificationSubscriberRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	subscriberID := d.Get("subscriber_id").(string)

	subscriber, resp, err := client.Subscriber.GetSubscriber(subscriberID) // Get all producers

	if err != nil {
		if resp == nil || resp.StatusCode != http.StatusForbidden { // Do not error on permission issues
			return diag.FromErr(err)
		}
		subscriber = &notification.Subscriber{}
	}
	d.SetId(subscriberID)
	_ = d.Set("managing_organization_id", subscriber.ManagingOrganizationID)
	_ = d.Set("managing_organization", subscriber.ManagingOrganization)
	_ = d.Set("subscriber_product_name", subscriber.SubscriberProductName)
	_ = d.Set("subscriber_service_name", subscriber.SubscriberServicename)
	_ = d.Set("subscriber_service_instance_name", subscriber.SubscriberServiceinstanceName)
	_ = d.Set("subscriber_service_base_url", subscriber.SubscriberServiceBaseURL)
	_ = d.Set("subscriber_service_path_url", subscriber.SubscriberServicePathURL)
	_ = d.Set("description", subscriber.Description)

	return diags
}
