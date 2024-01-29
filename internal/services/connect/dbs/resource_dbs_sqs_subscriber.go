package dbs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/philips-software/go-hsdp-api/connect/dbs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceDBSSQSSubscriber() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				_ = d.Set("name_infix", "imported")
				return []*schema.ResourceData{d}, nil
			},
		},
		CreateContext: resourceDBSSQSSubscriberCreate,
		ReadContext:   resourceDBSSQSSubscriberRead,
		DeleteContext: resourceDBSSQSSubscriberDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name_infix": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressWhenImported,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"queue_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delivery_delay_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ForceNew: true,
			},
			"message_retention_period_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  345600,
				ForceNew: true,
			},
			"receive_wait_time_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ForceNew: true,
			},
			"server_side_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"queue_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaToDBSSQSSubscriber(d *schema.ResourceData) dbs.SQSSubscriberConfig {
	nameInfix := d.Get("name_infix").(string)
	description := d.Get("description").(string)
	queueType := d.Get("queue_type").(string)
	deliveryDelaySeconds := d.Get("delivery_delay_seconds").(int)
	messageRetentionPeriod := d.Get("message_retention_period_seconds").(int)
	receiveWaitTimeSeconds := d.Get("receive_wait_time_seconds").(int)
	serverSideEncryption := d.Get("server_side_encryption").(bool)
	resource := dbs.SQSSubscriberConfig{
		ResourceType:                  "SQSSubscriberConfig",
		NameInfix:                     nameInfix,
		Description:                   description,
		QueueType:                     queueType,
		DeliveryDelaySeconds:          deliveryDelaySeconds,
		MessageRetentionPeriod:        messageRetentionPeriod,
		ReceiveMessageWaitTimeSeconds: receiveWaitTimeSeconds,
		ServerSideEncryption:          serverSideEncryption,
	}
	return resource
}

func dbsSQSSubscriberToSchema(resource dbs.SQSSubscriber, d *schema.ResourceData) {
	_ = d.Set("name", resource.Name)
	_ = d.Set("description", resource.Description)
	_ = d.Set("queue_type", resource.QueueType)
	_ = d.Set("delivery_delay_seconds", resource.DeliveryDelaySeconds)
	_ = d.Set("message_retention_period_seconds", resource.MessageRetentionPeriod)
	_ = d.Set("receive_wait_time_seconds", resource.ReceiveMessageWaitTimeSeconds)
	_ = d.Set("server_side_encryption", resource.ServerSideEncryption)
	_ = d.Set("status ", resource.Status)
	_ = d.Set("queue_name", resource.QueueName)
}

func resourceDBSSQSSubscriberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics
	principal := config.SchemaToPrincipal(d, m)

	client, err := c.DBSClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToDBSSQSSubscriber(d)

	var created *dbs.SQSSubscriber
	var resp *dbs.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.Subscribers.CreateSQS(resource)
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
	if created == nil {
		return diag.FromErr(fmt.Errorf("failed to create resource: %d", resp.StatusCode()))
	}
	d.SetId(created.ID)

	dbsSQSSubscriberToSchema(*created, d)
	return diags
}

func resourceDBSSQSSubscriberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.DBSClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	var resource *dbs.SQSSubscriber
	var resp *dbs.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.Subscribers.GetSQSByID(d.Id())
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp != nil && (resp.StatusCode() == http.StatusNotFound || resp.StatusCode() == http.StatusGone) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	dbsSQSSubscriberToSchema(*resource, d)
	return diags
}

func resourceDBSSQSSubscriberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.DBSClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	var resource *dbs.SQSSubscriber
	var resp *dbs.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.Subscribers.GetSQSByID(d.Id())
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

	var ok bool
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		ok, _, err = client.Subscribers.DeleteSQS(*resource)
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
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
