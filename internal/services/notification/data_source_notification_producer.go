package notification

import (
	"context"
	"net/http"

	"github.com/philips-software/go-dip-api/notification"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceNotificationProducer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationProducerRead,
		Schema: map[string]*schema.Schema{
			"producer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"managing_organization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"managing_organization": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"producer_product_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"producer_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"producer_service_instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"producer_service_base_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"producer_service_path_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceNotificationProducerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.NotificationClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	producerID := d.Get("producer_id").(string)

	producer, resp, err := client.Producer.GetProducer(producerID)

	if err != nil {
		if resp == nil || resp.StatusCode != http.StatusForbidden { // Do not error on permission issues
			return diag.FromErr(err)
		}
		producer = &notification.Producer{}
	}

	d.SetId(producerID)
	_ = d.Set("managing_organization_id", producer.ManagingOrganizationID)
	_ = d.Set("managing_organization", producer.ManagingOrganization)
	_ = d.Set("producer_product_name", producer.ProducerProductName)
	_ = d.Set("producer_service_name", producer.ProducerServiceName)
	_ = d.Set("producer_service_instance_name", producer.ProducerServiceInstanceName)
	_ = d.Set("producer_service_base_url", producer.ProducerServiceBaseURL)
	_ = d.Set("producer_service_path_url", producer.ProducerServicePathURL)
	_ = d.Set("description", producer.Description)

	return diags
}
