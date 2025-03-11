package mdm

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dip-software/go-dip-api/connect/mdm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceConnectMDMBlobSubscription() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceConnectMDMBlobSubscriptionCreate,
		ReadContext:   resourceConnectMDMBlobSubscriptionRead,
		UpdateContext: resourceConnectMDMBlobSubscriptionUpdate,
		DeleteContext: resourceConnectMDMBlobSubscriptionDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_type_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notification_topic_id": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: tools.SuppressMulti(
					tools.SuppressDefaultSystemValue),
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaToBlobSubscription(d *schema.ResourceData) mdm.BlobSubscription {
	name := d.Get("name").(string)
	dataTypeId := d.Get("data_type_id").(string)
	notificationTopicId := d.Get("notification_topic_id").(string)
	description := d.Get("description").(string)

	resource := mdm.BlobSubscription{
		Name:        name,
		Description: description,
		DataTypeId:  mdm.Reference{Reference: dataTypeId},
	}
	if len(notificationTopicId) > 0 {
		identifier := mdm.Identifier{}
		parts := strings.Split(notificationTopicId, "|")
		if len(parts) > 1 {
			identifier.System = parts[0]
			identifier.Value = parts[1]
		} else {
			identifier.Value = notificationTopicId
		}
		resource.NotificationTopicGuid = identifier
	}
	return resource
}

func blobSubscriptionToSchema(resource mdm.BlobSubscription, d *schema.ResourceData) {
	_ = d.Set("description", resource.Description)
	_ = d.Set("name", resource.Name)
	_ = d.Set("data_type_id", resource.DataTypeId.Reference)
	_ = d.Set("guid", resource.ID)
	if resource.NotificationTopicGuid.Value != "" {
		value := resource.NotificationTopicGuid.Value
		if resource.NotificationTopicGuid.System != "" {
			value = fmt.Sprintf("%s|%s", resource.NotificationTopicGuid.System, resource.NotificationTopicGuid.Value)
		}
		_ = d.Set("notification_topic_id", value)
	}
}

func resourceConnectMDMBlobSubscriptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToBlobSubscription(d)

	var created *mdm.BlobSubscription
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.BlobSubscriptions.Create(resource)
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
	d.SetId(fmt.Sprintf("BlobSubscription/%s", created.ID))
	return resourceConnectMDMBlobSubscriptionRead(ctx, d, m)
}

func resourceConnectMDMBlobSubscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "BlobSubscription/%s", &id)
	var resource *mdm.BlobSubscription
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.BlobSubscriptions.GetByID(id)
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
	blobSubscriptionToSchema(*resource, d)
	return diags
}

func resourceConnectMDMBlobSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)

	resource := schemaToBlobSubscription(d)
	resource.ID = id

	_, _, err = client.BlobSubscriptions.Update(resource)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		return diags
	}
	return resourceConnectMDMBlobSubscriptionRead(ctx, d, m)
}

func resourceConnectMDMBlobSubscriptionDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	resource, _, err := client.BlobSubscriptions.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.BlobSubscriptions.Delete(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
