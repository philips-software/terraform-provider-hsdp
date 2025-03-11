package notification

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dip-software/go-dip-api/notification"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceNotificationTopic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationTopicRead,
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
	var topic *notification.Topic
	var resp *notification.Response
	var err error

	client, err := c.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	topicID := d.Get("topic_id").(string)
	name := d.Get("name").(string)

	if topicID == "" && name == "" {
		return diag.FromErr(fmt.Errorf("need either a topic_id or a name to query"))
	}
	if topicID != "" && name != "" {
		return diag.FromErr(fmt.Errorf("specify either a topic_id or a name, not both"))
	}

	if topicID != "" {
		topic, resp, err = client.Topic.GetTopic(topicID) // Get topic
		if err != nil {
			if resp == nil || resp.StatusCode != http.StatusForbidden { // Do not error on permission issues
				return diag.FromErr(err)
			}
			return diag.FromErr(err)
		}
	}
	if name != "" {
		opts := &notification.GetOptions{
			Name: &name,
		}
		list, resp, err := client.Topic.GetTopics(opts) // Get all topics
		if err != nil {
			if resp == nil || resp.StatusCode != http.StatusForbidden { // Do not error on permission issues
				return diag.FromErr(err)
			}
			list = []notification.Topic{} // empty list
		}
		if len(list) == 0 {
			return diag.FromErr(fmt.Errorf("no matching topic found"))
		}
		topic, resp, err = client.Topic.GetTopic(list[0].ID) // Get topic
		if err != nil {
			if resp == nil || resp.StatusCode != http.StatusForbidden { // Do not error on permission issues
				return diag.FromErr(err)
			}
			return diag.FromErr(err)
		}
	}
	d.SetId(topic.ID)
	_ = d.Set("name", topic.Name)
	_ = d.Set("producer_id", topic.ProducerID)
	_ = d.Set("scope", topic.Scope)
	_ = d.Set("allowed_scopes", topic.AllowedScopes)
	_ = d.Set("is_auditable", topic.IsAuditable)
	_ = d.Set("description", topic.Description)

	return diags
}
