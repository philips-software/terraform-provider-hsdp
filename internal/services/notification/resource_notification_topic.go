package notification

import (
	"context"
	"errors"
	"net/http"

	"github.com/dip-software/go-dip-api/notification"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceNotificationTopic() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
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
			"principal": config.PrincipalSchema(),
			"is_auditable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"soft_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNotificationTopicDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	principal := config.SchemaToPrincipal(d, m)

	client, err := c.NotificationClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	_, resp, err := client.Topic.DeleteTopic(notification.Topic{ID: d.Id()})
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

func resourceNotificationTopicRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	principal := config.SchemaToPrincipal(d, m)

	client, err := c.NotificationClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	topic, _, err := client.Topic.GetTopic(d.Id())
	if err != nil {
		if errors.Is(err, notification.ErrEmptyResult) { // Removed
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
	c := m.(*config.Config)
	principal := config.SchemaToPrincipal(d, m)

	client, err := c.NotificationClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	topic := notification.Topic{
		Name:          d.Get("name").(string),
		Scope:         d.Get("scope").(string),
		ProducerID:    d.Get("producer_id").(string),
		AllowedScopes: tools.ExpandStringList(d.Get("allowed_scopes").(*schema.Set).List()),
		IsAuditable:   d.Get("is_auditable").(bool),
		Description:   d.Get("description").(string),
	}

	var created *notification.Topic

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var resp *notification.Response

		created, resp, err = client.Topic.CreateTopic(topic)
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
	return resourceNotificationTopicRead(ctx, d, m)
}

func resourceNotificationTopicUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	principal := config.SchemaToPrincipal(d, m)

	client, err := c.NotificationClient(principal)
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
		topic.AllowedScopes = tools.ExpandStringList(d.Get("allowed_scopes").(*schema.Set).List())
		topic.Description = d.Get("description").(string)
		topic.IsAuditable = d.Get("is_auditable").(bool)
		err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
			var resp *notification.Response
			_, resp, err = client.Topic.UpdateTopic(*topic)
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
	}
	return diags
}
