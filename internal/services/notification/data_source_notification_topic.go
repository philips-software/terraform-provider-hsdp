package notification

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/notification"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceNotificationTopic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationTopicRead,
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"producer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allowed_scopes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_auditable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNotificationTopicRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	topicID := d.Get("topic_id").(string)

	topic, resp, err := client.Topic.GetTopic(topicID) // Get all producers

	if err != nil {
		if resp == nil || resp.StatusCode != http.StatusForbidden { // Do not error on permission issues
			return diag.FromErr(err)
		}
		topic = &notification.Topic{}
	}
	d.SetId(topicID)
	_ = d.Set("name", topic.Name)
	_ = d.Set("producer_id", topic.ProducerID)
	_ = d.Set("scope", topic.Scope)
	_ = d.Set("allowed_scopes", topic.AllowedScopes)
	_ = d.Set("is_auditable", topic.IsAuditable)
	_ = d.Set("description", topic.Description)

	return diags
}
