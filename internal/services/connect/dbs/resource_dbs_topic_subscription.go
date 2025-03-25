package dbs

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"net/http"
	"regexp"
	"time"

	"github.com/dip-software/go-dip-api/connect/dbs"

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
		CreateContext:      resourceDBSTopicSubscriptionCreate,
		ReadContext:        resourceDBSTopicSubscriptionRead,
		DeleteContext:      resourceDBSTopicSubscriptionDelete,
		DeprecationMessage: "This resource is deprecated. It will be removed in an upcoming release.",
		SchemaVersion:      1,
		Schema: map[string]*schema.Schema{
			"name_infix": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressWhenImported,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 12),
					validation.StringMatch(regexp.MustCompile("^[a-zA-Z0-9]+$"), "value must be alphanumeric"),
				),
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 250),
					validation.StringMatch(regexp.MustCompile("^[-a-zA-Z0-9_, .]+$"), ""),
				),
			},
			"subscriber_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 36),
					validation.StringMatch(regexp.MustCompile("^[-a-fA-F0-9]+$"), ""),
				),
			},
			"data_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 36),
					validation.StringMatch(regexp.MustCompile("^[-a-zA-Z0-9_.]+$"), ""),
				),
			},
			"deliver_data_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"kinesis_stream_partition_key": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 250),
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
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.Subscriptions.CreateTopicSubscription(resource)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	}, append(tools.StandardRetryOnCodes, http.StatusUnprocessableEntity)...)
	if err != nil {
		return diag.FromErr(err)
	}
	if created == nil {
		return diag.FromErr(fmt.Errorf("failed to create resource: %d", resp.StatusCode()))
	}
	d.SetId(created.ID)

	created, err = waitResourceCreated[dbs.TopicSubscription](ctx, StatusTopicSubscription(ctx, client, d.Id()),
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		if created != nil {
			return diag.FromErr(fmt.Errorf("resource did not get correct state: %s", created.ErrorMessage))
		}
		return diag.FromErr(err)
	}

	dbsTopicSubscriptionToSchema(*created, d)
	return diags
}

func waitResourceCreated[V dbs.TopicSubscription | dbs.SQSSubscriber](
	ctx context.Context,
	refresh retry.StateRefreshFunc,
	timeout time.Duration) (*V, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"Creating", "Updating"},
		Target:     []string{"Active"},
		Refresh:    refresh,
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}
	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*V); ok {
		return output, err
	}

	return nil, err
}

func StatusTopicSubscription(ctx context.Context, client *dbs.Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resource *dbs.TopicSubscription
		var resp *dbs.Response
		err := tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
			var err error
			resource, resp, err = client.Subscriptions.GetTopicSubscriptionByID(id)
			if err != nil {
				_ = client.TokenRefresh()
			}
			if resp == nil {
				return nil, err
			}
			return resp.Response, err
		})

		if err != nil {
			return nil, "", err
		}

		return resource, resource.Status, nil
	}
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

func resourceDBSTopicSubscriptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	var ok bool
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		ok, _, err = client.Subscriptions.DeleteTopicSubscription(*resource)
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
