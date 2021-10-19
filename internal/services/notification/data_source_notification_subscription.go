package notification

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/notification"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceNotificationSubscription() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationSubscriptionRead,
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"topic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscriber_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscription_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNotificationSubscriptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	subscriptionID := d.Get("subscription_id").(string)

	subscription, resp, err := client.Subscription.GetSubscription(subscriptionID)

	if err != nil {
		if resp == nil || resp.StatusCode != http.StatusForbidden { // Do not error on permission issues
			return diag.FromErr(err)
		}
		subscription = &notification.Subscription{}
	}
	d.SetId(subscriptionID)
	_ = d.Set("topic_id", subscription.TopicID)
	_ = d.Set("subscriber_id", subscription.SubscriberID)
	_ = d.Set("subscription_endpoint", subscription.SubscriptionEndpoint)

	return diags
}
