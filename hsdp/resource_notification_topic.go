package hsdp

import (
	"context"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/notification"
)

func resourceNotificationTopic() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceNotificationTopicCreate,
		ReadContext:   resourceNotificationTopicRead,
		UpdateContext: resourceNotificationTopicUpdate,
		DeleteContext: resourceNotificationTopicDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"producer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"allowed_scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_auditable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceNotificationTopicDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	client, err := config.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	_, _, err = client.Topic.DeleteTopic(notification.Topic{ID: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func resourceNotificationTopicRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	client, err := config.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	topic, _, err := client.Topic.GetTopic(d.Id())
	if err != nil {
		if err == notification.ErrEmptyResult { // Removed
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("name", topic.Name)
	_ = d.Set("producer_id", topic.ProducerID)
	_ = d.Set("scope", topic.Scope)
	_ = d.Set("allowed_scopes", topic.AllowedScopes)
	_ = d.Set("is_auditable", topic.IsAuditable)
	_ = d.Set("description", topic.Description)

	return diags
}

func resourceNotificationTopicCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client, err := config.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	topic := notification.Topic{
		Name:          d.Get("name").(string),
		Scope:         d.Get("scope").(string),
		ProducerID:    d.Get("producer_id").(string),
		AllowedScopes: expandStringList(d.Get("allowed_scopes").(*schema.Set).List()),
		IsAuditable:   d.Get("is_auditable").(bool),
		Description:   d.Get("description").(string),
	}

	var created *notification.Topic

	operation := func() error {
		var resp *notification.Response
		_ = client.TokenRefresh()
		created, resp, err = client.Topic.CreateTopic(topic)
		return checkForNotificationPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(created.ID)
	return resourceNotificationTopicRead(ctx, d, m)
}

func resourceNotificationTopicUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	client, err := config.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	topic, _, err := client.Topic.GetTopic(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("name") || d.HasChange("allowed_scopes") || d.HasChange("description") ||
		d.HasChange("is_auditable") {
		topic.Name = d.Get("name").(string)
		topic.AllowedScopes = expandStringList(d.Get("allowed_scopes").(*schema.Set).List())
		topic.Description = d.Get("description").(string)
		topic.IsAuditable = d.Get("is_auditable").(bool)
		operation := func() error {
			var resp *notification.Response
			_ = client.TokenRefresh()
			_, _, err = client.Topic.UpdateTopic(*topic)
			return checkForNotificationPermissionErrors(client, resp, err)
		}
		err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}
