package notification

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dip-software/go-dip-api/notification"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceNotificationSubscription() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceNotificationSubscriptionCreate,
		UpdateContext: resourceNotificationSubscriptionUpdate,
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
			"principal": config.PrincipalSchema(),
			"subscription_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"soft_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNotificationSubscriptionUpdate(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if !d.HasChanges("soft_delete") {
		return diag.FromErr(fmt.Errorf("only 'soft_delete' can be updated, this is provider bug"))
	}
	return diags
}

func resourceNotificationSubscriptionDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	principal := config.SchemaToPrincipal(d, m)

	client, err := c.NotificationClient(principal)
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

func resourceNotificationSubscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	principal := config.SchemaToPrincipal(d, m)

	client, err := c.NotificationClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	var subscription *notification.Subscription

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var resp *notification.Response
		subscription, resp, err = client.Subscription.GetSubscription(d.Id())
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, fmt.Errorf("Producer.GetSubscription: response is nil, error: %w", err)
		}
		return resp.Response, err
	})
	if err != nil {
		if errors.Is(err, notification.ErrEmptyResult) { // Removed
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
	c := m.(*config.Config)
	principal := config.SchemaToPrincipal(d, m)

	client, err := c.NotificationClient(principal)
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

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var resp *notification.Response
		created, resp, err = client.Subscription.CreateSubscription(subscription)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, fmt.Errorf("Subscription.CreateSubscription: response is nil")
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(created.ID)
	return resourceNotificationSubscriptionRead(ctx, d, m)
}
