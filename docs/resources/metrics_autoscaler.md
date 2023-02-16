---
subcategory: "Autoscaler"
page_title: "HSDP: hsdp_metrics_autoscaler"
description: |-
  Manages HSDP Metrics autoscalers
---

# hsdp_metrics_autoscaler

Manages HSDP Metrics Autoscaler settings for Cloudfoundry applications hosted in an HSDP CF space.

~> **NOTE:** This resource is only available when the `region` and `uaa_*` keys are set in the provider config

[Metrics Service Broker](https://www.hsdp.io/documentation/metrics-service-broker)

## Example Usage

The following resource enables autoscaling of the HSDP CF hosted `myapp`, scaling up to 10 instances maximum.
The app upscales at 90% CPU utilization and downscales again when CPU usage falls below 5%.

```hcl
resource "hsdp_metrics_autoscaler" "myapp_autoscaler" {
  metrics_instance_id = cloudfoundry_service_instance.metrics.id
  app_name    = cloudfoundry_app.myapp.name
 
  enabled = true
  min_instances     = 1
  max_instances     = 10 

  threshold_cpu {
    enabled = true
    min     = 5
    max     = 90
  }

  threshold_memory {
    enabled = false
  }

  threshold_http_latency {
    enabled = false
  }

  threshold_http_rate {
    enabled = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `metrics_instance_id` - (Required) The Metrics service instance UUID running in the space where the app is hosted.
* `app_name` - (Required) The CF app name to apply this autoscaler settings for.
* `min` - (Optional) Minimum number of app instances. Default: 1
* `max` - (Optional) Maximum number of app instances. Default: 10
* `threshold_cpu` - (Required) CPU threshold block. Min/max values are `%`
* `threshold_memory` - (Required) Memory threshold block. Min/max values are `%`
* `threshold_http_latency` - (Required) HTTP latency threshold block. Min/max values are in `ms`
* `threshold_http_rate` - (Required) HTTP rate threshold block. Min/max values are in `requests/minute`

For each threshold block the following argments are supported:

* `enabled` - (Required) When set to `true` this threshold type is evaluated
* `min` - (Optional) The minimum value of the resource. When the resource hits this value, downscaling is triggered.
* `max` - (Optional) The maxmium value of the resource. When the resource hits this value, upscaling is triggered.

## Attributes Reference

The following attributes are exported:

* `id` - The resource instance ID

## Import

Instance state is imported automatically
