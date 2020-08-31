package hsdp

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/philips-software/go-hsdp-api/console"
	"time"
)

func resourceMetricsAutoscaler() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceMetricsAutoscalerCreate,
		Read:   resourceMetricsAutoscalerRead,
		Update: resourceMetricsAutoscalerUpdate,
		Delete: resourceMetricsAutoscalerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(12 * time.Minute),
			Update: schema.DefaultTimeout(12 * time.Minute),
			Delete: schema.DefaultTimeout(22 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"max_instances": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},
			"min_instances": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"thresholds": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"max": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  0,
						},
						"min": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  0,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceMetricsAutoscalerDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.ConsoleClient()
	if err != nil {
		return err
	}
	var app console.Application
	instanceID := d.Get("instance_id").(string)
	app.Name = d.Get("name").(string)
	app.Enabled = false
	result, _, err := client.Metrics.UpdateApplicationAutoscaler(instanceID, app)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("error creating/updating autoscaler")
	}
	d.SetId("")
	return nil
}

func resourceMetricsAutoscalerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceMetricsAutoscalerCreate(d, m)
}

func resourceMetricsAutoscalerRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.ConsoleClient()
	if err != nil {
		return err
	}
	instanceID := d.Get("instance_id").(string)
	name := d.Get("name").(string)

	app, _, err := client.Metrics.GetApplicationAutoscaler(instanceID, name)
	_ = d.Set("min_instances", app.MinInstances)
	_ = d.Set("max_instances", app.MaxInstances)
	_ = d.Set("enabled", app.Enabled)

	return nil
}

func resourceMetricsAutoscalerCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.ConsoleClient()
	if err != nil {
		return err
	}
	var app console.Application

	instanceID := d.Get("instance_id").(string)

	app.Name = d.Get("name").(string)
	app.MaxInstances = d.Get("max_instances").(int)
	app.MinInstances = d.Get("min_instances").(int)
	app.Enabled = d.Get("enabled").(bool)

	created, _, err := client.Metrics.UpdateApplicationAutoscaler(instanceID, app)
	if err != nil {
		return err
	}
	if created == nil {
		return fmt.Errorf("error creating/updating autoscaler")
	}
	d.SetId(instanceID + created.Name)
	return nil
}
