package hsdp

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/notification"
)

func dataSourceNotificationTopics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationTopicsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"topic_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceNotificationTopicsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var diags diag.Diagnostics

	client, err := config.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	topicName := d.Get("name").(string)

	opts := &notification.GetOptions{
		Name: &topicName,
	}

	list, resp, err := client.Topic.GetTopics(opts) // Get all producers

	if err != nil {
		if resp == nil || resp.StatusCode != http.StatusForbidden { // Do not error on permission issues
			return diag.FromErr(err)
		}
		list = []notification.Topic{} // empty list
	}
	topicIDs := make([]string, 0)

	for _, p := range list {
		topicIDs = append(topicIDs, p.ID)
	}
	d.SetId("notification_topics")
	_ = d.Set("topic_ids", topicIDs)

	return diags
}
