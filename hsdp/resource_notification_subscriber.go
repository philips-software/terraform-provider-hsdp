package hsdp

import (
	"context"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/notification"
)

func resourceNotificationSubscriber() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceNotificationSubscriberCreate,
		ReadContext:   resourceNotificationSubscriberRead,
		DeleteContext: resourceNotificationSubscriberDelete,

		Schema: map[string]*schema.Schema{
			"managing_organization_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"managing_organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subscriber_product_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subscriber_service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subscriber_service_instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subscriber_service_base_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subscriber_service_path_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNotificationSubscriberDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	client, err := config.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	_, _, err = client.Subscriber.DeleteSubscriber(notification.Subscriber{ID: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func resourceNotificationSubscriberRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	client, err := config.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	var subscriber *notification.Subscriber

	operation := func() error {
		var resp *notification.Response
		var err error
		_ = client.TokenRefresh()
		subscriber, resp, err = client.Subscriber.GetSubscriber(d.Id())
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

func resourceNotificationSubscriberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client, err := config.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	subscriber := notification.Subscriber{
		ManagingOrganizationID:        d.Get("managing_organization_id").(string),
		ManagingOrganization:          d.Get("managing_organization").(string),
		SubscriberProductName:         d.Get("subscriber_product_name").(string),
		SubscriberServicename:         d.Get("subscriber_service_name").(string),
		SubscriberServiceinstanceName: d.Get("subscriber_service_instance_name").(string),
		SubscriberServiceBaseURL:      d.Get("subscriber_service_base_url").(string),
		SubscriberServicePathURL:      d.Get("subscriber_service_path_url").(string),
		Description:                   d.Get("description").(string),
	}

	var created *notification.Subscriber

	operation := func() error {
		var resp *notification.Response
		_ = client.TokenRefresh()
		created, resp, err = client.Subscriber.CreateSubscriber(subscriber)
		return checkForNotificationPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(created.ID)
	return resourceNotificationSubscriberRead(ctx, d, m)
}
