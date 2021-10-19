package notification

import (
	"context"
	"net/http"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/notification"
	config2 "github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceNotificationSubscription() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceNotificationSubscriptionCreate,
		ReadContext:   resourceNotificationSubscriptionRead,
		DeleteContext: resourceNotificationSubscriptionDelete,

		Schema: map[string]*schema.Schema{
			"topic_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subscriber_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subscription_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"soft_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
		},
	}
}

func resourceNotificationSubscriptionDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*config2.Config)
	client, err := config.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	_, resp, err := client.Subscription.DeleteSubscription(notification.Subscription{ID: d.Id()})
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusForbidden {
			softDelete := d.Get("soft_delete").(bool)
			if softDelete { // No error on delete
				d.SetId("")
				return diags
			}
		}
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func resourceNotificationSubscriptionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*config2.Config)
	client, err := config.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	var subscription *notification.Subscription

	operation := func() error {
		var resp *notification.Response
		var err error
		_ = client.TokenRefresh()
		subscription, resp, err = client.Subscription.GetSubscription(d.Id())
		return checkForNotificationPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		if err == notification.ErrEmptyResult { // Removed
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("topic_id", subscription.TopicID)
	_ = d.Set("subscriber_id", subscription.SubscriberID)
	_ = d.Set("subscription_endpoint", subscription.SubscriptionEndpoint)

	return diags
}

func resourceNotificationSubscriptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*config2.Config)
	client, err := config.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	subscription := notification.Subscription{
		TopicID:              d.Get("topic_id").(string),
		SubscriberID:         d.Get("subscriber_id").(string),
		SubscriptionEndpoint: d.Get("subscription_endpoint").(string),
	}

	var created *notification.Subscription

	operation := func() error {
		var resp *notification.Response
		_ = client.TokenRefresh()
		created, resp, err = client.Subscription.CreateSubscription(subscription)
		return checkForNotificationPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(created.ID)
	return resourceNotificationSubscriptionRead(ctx, d, m)
}
