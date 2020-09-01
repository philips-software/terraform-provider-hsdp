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
			"metrics_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"app_name": {
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
			"threshold_http_latency": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"max": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  10000,
						},
						"min": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  10,
						},
					},
				},
			},
			"threshold_http_rate": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"max": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  6000000,
						},
						"min": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  300,
						},
					},
				},
			},
			"threshold_memory": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"max": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  100,
						},
						"min": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  20,
						},
					},
				},
			},
			"threshold_cpu": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
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
							Default:  100,
						},
						"min": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  5,
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
	instanceID := d.Get("metrics_instance_id").(string)
	app.Name = d.Get("app_name").(string)
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
	instanceID := d.Get("metrics_instance_id").(string)
	name := d.Get("app_name").(string)

	app, _, err := client.Metrics.GetApplicationAutoscaler(instanceID, name)
	if err != nil {
		return err
	}
	_ = d.Set("min_instances", app.MinInstances)
	_ = d.Set("max_instances", app.MaxInstances)
	_ = d.Set("enabled", app.Enabled)
	for _, th := range app.Thresholds {
		mapping, ok := thresholdMapping[th.Name]
		if !ok {
			return fmt.Errorf("unknown threshold: %s", th.Name)
		}
		fields := make(map[string]interface{})
		fields["enabled"] = th.Enabled
		fields["min"] = th.Min
		fields["max"] = th.Max
		s := &schema.Set{F: resourceMetricsThresholdHash}
		s.Add(fields)
		_ = d.Set(mapping, s)
	}
	return nil
}

func resourceMetricsThresholdHash(v interface{}) int {
	return 0
}

func resourceMetricsAutoscalerCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.ConsoleClient()
	if err != nil {
		return err
	}
	var app console.Application

	instanceID := d.Get("metrics_instance_id").(string)

	app.Name = d.Get("app_name").(string)
	app.MaxInstances = d.Get("max_instances").(int)
	app.MinInstances = d.Get("min_instances").(int)
	app.Enabled = d.Get("enabled").(bool)

	for key, mapping := range thresholdMapping {
		if v, ok := d.GetOk(mapping); ok {
			vL := v.(*schema.Set).List()
			for _, vi := range vL {
				mVi := vi.(map[string]interface{})
				var threshold console.Threshold
				threshold.Name = key
				threshold.Min = mVi["min"].(float64)
				threshold.Max = mVi["max"].(float64)
				threshold.Enabled = mVi["enabled"].(bool)
				app.Thresholds = append(app.Thresholds, threshold)
			}
		}
	}
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

func expandMapList(configured []interface{}) []map[string]interface{} {
	vs := make([]map[string]interface{}, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, val)
		}
	}
	return vs
}

func stringInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
