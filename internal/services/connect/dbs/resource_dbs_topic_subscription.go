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
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceDBSTopicSubscriptionCreate,
		ReadContext:   resourceDBSTopicSubscriptionyRead,
		DeleteContext: resourceDBSTopicSubscriptionDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// TODO: fill in
		},
	}
}

func schemaToDBSTopicSubscription(d *schema.ResourceData) dbs.TopicSubscriptionConfig {

	resource := dbs.TopicSubscriptionConfig{
		ResourceType: "Subscription",
		// TODO: fill in
	}
	return resource
}

func dbsTopicSubscriptionToSchema(resource dbs.TopicSubscription, d *schema.ResourceData) {
	// TODO: fill in
}

func resourceDBSTopicSubscriptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

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
	_ = d.Set("guid", created.ID)
	d.SetId(fmt.Sprintf("BlobStorePolicy/%s", created.ID))

	return resourceDBSTopicSubscriptionyRead(ctx, d, m)
}

func resourceDBSTopicSubscriptionyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.DBSClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "BlobStorePolicy/%s", &id)
	var resource *dbs.TopicSubscription
	var resp *dbs.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
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

	var id string
	_, _ = fmt.Sscanf(d.Id(), "BlobStorePolicy/%s", &id)
	resource, _, err := client.Subscriptions.GetTopicSubscriptionByID(id)
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
