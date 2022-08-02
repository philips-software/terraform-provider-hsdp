package notification

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/notification"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceNotificationProducer() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceNotificationProducerCreate,
		ReadContext:   resourceNotificationProducerRead,
		DeleteContext: resourceNotificationProducerDelete,

		Schema: map[string]*schema.Schema{
			"principal": config.PrincipalSchema(),
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
			"producer_product_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"producer_service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"producer_service_instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"producer_service_base_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"producer_service_path_url": {
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

func resourceNotificationProducerDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.NotificationClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	_, resp, err := client.Producer.DeleteProducer(notification.Producer{ID: d.Id()})
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

func resourceNotificationProducerRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.NotificationClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	var producer *notification.Producer

	operation := func() error {
		var resp *notification.Response
		_ = client.TokenRefresh()
		producer, resp, err = client.Producer.GetProducer(d.Id())
		return checkForNotificationPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10))
	if err != nil {
		if errors.Is(err, notification.ErrEmptyResult) { // Removed
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("managing_organization_id", producer.ManagingOrganizationID)
	_ = d.Set("managing_organization", producer.ManagingOrganization)
	_ = d.Set("producer_product_name", producer.ProducerProductName)
	_ = d.Set("producer_service_name", producer.ProducerServiceName)
	_ = d.Set("producer_service_instance_name", producer.ProducerServiceInstanceName)
	_ = d.Set("producer_service_base_url", producer.ProducerServiceBaseURL)
	_ = d.Set("producer_service_path_url", producer.ProducerServicePathURL)
	_ = d.Set("description", producer.Description)

	return diags
}

func resourceNotificationProducerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	principal := config.SchemaToPrincipal(d, m)

	client, err := c.NotificationClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	producer := notification.Producer{
		ManagingOrganizationID:      d.Get("managing_organization_id").(string),
		ManagingOrganization:        d.Get("managing_organization").(string),
		ProducerProductName:         d.Get("producer_product_name").(string),
		ProducerServiceName:         d.Get("producer_service_name").(string),
		ProducerServiceInstanceName: d.Get("producer_service_instance_name").(string),
		ProducerServiceBaseURL:      d.Get("producer_service_base_url").(string),
		ProducerServicePathURL:      d.Get("producer_service_path_url").(string),
		Description:                 d.Get("description").(string),
	}

	var created *notification.Producer

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var resp *notification.Response
		created, resp, err = client.Producer.CreateProducer(producer)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, fmt.Errorf("Producer.CreateProducer: response is nil")
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(created.ID)
	return resourceNotificationProducerRead(ctx, d, m)
}

func checkForNotificationPermissionErrors(client *notification.Client, resp *notification.Response, err error) error {
	if resp == nil || resp.StatusCode > 500 {
		return err
	}
	if resp.StatusCode == http.StatusForbidden {
		_ = client.TokenRefresh()
		return err
	}
	return backoff.Permanent(err)
}
