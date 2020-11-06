package hsdp

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/console"
	"time"
)

func resourceMetricsAutoscaler() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMetricsAutoscalerCreate,
		ReadContext:   resourceMetricsAutoscalerRead,
		UpdateContext: resourceMetricsAutoscalerUpdate,
		DeleteContext: resourceMetricsAutoscalerDelete,

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

func resourceMetricsAutoscalerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.ConsoleClient()
	if err != nil {
		return diag.FromErr(err)
	}
	var app console.Application
	instanceID := d.Get("metrics_instance_id").(string)
	app.Name = d.Get("app_name").(string)
	app.Enabled = false
	result, _, err := client.Metrics.UpdateApplicationAutoscaler(instanceID, app)
	if err != nil {
		return diag.FromErr(err)
	}
	if result == nil {
		return diag.FromErr(fmt.Errorf("error creating/updating autoscaler"))
	}
	d.SetId("")
	return diags
}

func resourceMetricsAutoscalerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMetricsAutoscalerCreate(ctx, d, m)
}

func resourceMetricsAutoscalerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.ConsoleClient()
	if err != nil {
		return diag.FromErr(err)
	}
	instanceID := d.Get("metrics_instance_id").(string)
	name := d.Get("app_name").(string)

	app, _, err := client.Metrics.GetApplicationAutoscaler(instanceID, name)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("min_instances", app.MinInstances)
	_ = d.Set("max_instances", app.MaxInstances)
	_ = d.Set("enabled", app.Enabled)
	for _, th := range app.Thresholds {
		mapping, ok := thresholdMapping[th.Name]
		if !ok {
			return diag.FromErr(fmt.Errorf("unknown threshold: %s", th.Name))
		}
		fields := make(map[string]interface{})
		fields["enabled"] = th.Enabled
		fields["min"] = th.Min
		fields["max"] = th.Max
		s := &schema.Set{F: resourceMetricsThresholdHash}
		s.Add(fields)
		_ = d.Set(mapping, s)
	}
	return diags
}

func resourceMetricsThresholdHash(v interface{}) int {
	return 0
}

func resourceMetricsAutoscalerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.ConsoleClient()
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
	}
	if created == nil {
		return diag.FromErr(fmt.Errorf("error creating/updating autoscaler"))
	}
	d.SetId(instanceID + created.Name)
	return diags
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
