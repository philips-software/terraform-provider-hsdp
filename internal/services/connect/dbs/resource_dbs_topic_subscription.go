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

func ResourceDBSTopicSubscription() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				_ = d.Set("name_infix", "imported")
				return []*schema.ResourceData{d}, nil
			},
		},
		CreateContext: resourceDBSTopicSubscriptionCreate,
		ReadContext:   resourceDBSTopicSubscriptionRead,
		DeleteContext: resourceDBSTopicSubscriptionDelete,
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
			"subscriber_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deliver_data_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"kinesis_stream_partition_key": {
				Type:     schema.TypeString,
				Optional: true,
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
			"rule_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaToDBSTopicSubscription(d *schema.ResourceData) dbs.TopicSubscriptionConfig {
	nameInfix := d.Get("name_infix").(string)
	description := d.Get("description").(string)
	subscriberId := d.Get("subscriber_id").(string)
	dataType := d.Get("data_type").(string)
	deliverDataOnly := d.Get("deliver_data_only").(bool)
	kinesisStreamPartitionKey := d.Get("kinesis_stream_partition_key").(string)
	resource := dbs.TopicSubscriptionConfig{
		ResourceType:              "Subscription",
		NameInfix:                 nameInfix,
		Description:               description,
		SubscriberId:              subscriberId,
		DeliverDataOnly:           deliverDataOnly,
		KinesisStreamPartitionKey: kinesisStreamPartitionKey,
		DataType:                  dataType,
	}
	return resource
}

func dbsTopicSubscriptionToSchema(resource dbs.TopicSubscription, d *schema.ResourceData) {
	_ = d.Set("name", resource.Name)
	_ = d.Set("description", resource.Description)
	_ = d.Set("subscriber_id", resource.Subscriber.ID)
	_ = d.Set("deliver_data_only", resource.DeliverDataOnly)
	_ = d.Set("kinesis_stream_partition_key", resource.KinesisStreamPartitionKey)
	_ = d.Set("status", resource.Status)
	_ = d.Set("data_type", resource.DataType)
	_ = d.Set("rule_name", resource.RuleName)
}

func resourceDBSTopicSubscriptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics
	principal := config.SchemaToPrincipal(d, m)

	client, err := c.DBSClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToDBSTopicSubscription(d)

	var created *dbs.TopicSubscription
	var resp *dbs.Response
	err = tools.TryHTTPCall(ctx, 5, func() (*http.Response, error) {
		var err error
		created, resp, err = client.Subscriptions.CreateTopicSubscription(resource)
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

	dbsTopicSubscriptionToSchema(*created, d)
	return diags
}

func resourceDBSTopicSubscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.DBSClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	var resource *dbs.TopicSubscription
	var resp *dbs.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.Subscriptions.GetTopicSubscriptionByID(d.Id())
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
	dbsTopicSubscriptionToSchema(*resource, d)
	return diags
}

func resourceDBSTopicSubscriptionDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.DBSClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	resource, _, err := client.Subscriptions.GetTopicSubscriptionByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.Subscriptions.DeleteTopicSubscription(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
