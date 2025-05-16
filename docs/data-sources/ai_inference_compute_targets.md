---
subcategory: "AI Inference"
---

# hsdp_ai_inference_compute_targets

-> **Deprecation Notice** This data source is deprecated and will be removed in an upcoming release of the provider

Retrieves AI Inference Compute Targets

## Example usage

```hcl
data "hsdp_config" "inference" {
  service = "inference"
}

data "hsdp_ai_inference_service_instance" "inference" {
  base_url        = data.hsdp_config.inference.url
  organization_id = var.inference_tenant_org_id
}

data "hsdp_ai_inference_compute_targets" "targets" {
  endpoint = data.hsdp_ai_inference_service_instance.inference.endpoint
}
```

## Argument reference

* `endpoint`- (Required) the AI Inference endpoint

## Attribute reference

The following attributes are exported:

* `ids` -  The list of container host IDs
* `names` - The list of container host names. This matches up with the `ids` list index.
