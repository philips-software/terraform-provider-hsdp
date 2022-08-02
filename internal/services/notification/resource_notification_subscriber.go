package notification

import (
	"context"
	"errors"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/notification"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceNotificationSubscriber() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
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
			"principal": config.PrincipalSchema(),
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
			"soft_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
		},
	}
}

func resourceNotificationSubscriberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	principal := config.SchemaToPrincipal(d, m)

	client, err := c.NotificationClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	var resp *notification.Response

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		_, resp, err = client.Subscriber.DeleteSubscriber(notification.Subscriber{ID: d.Id()})
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusForbidden {
			softDelete := d.Get("soft_delete").(bool)
			if softDelete {
				d.SetId("")
				return diags
			}
		}
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func resourceNotificationSubscriberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	principal := config.SchemaToPrincipal(d, m)

	client, err := c.NotificationClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	var subscriber *notification.Subscriber

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var resp *notification.Response
		subscriber, resp, err = client.Subscriber.GetSubscriber(d.Id())
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
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
	c := m.(*config.Config)
	principal := config.SchemaToPrincipal(d, m)

	client, err := c.NotificationClient(principal)
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

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var resp *notification.Response
		created, resp, err = client.Subscriber.CreateSubscriber(subscriber)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(created.ID)
	return resourceNotificationSubscriberRead(ctx, d, m)
}
