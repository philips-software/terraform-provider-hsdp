---
subcategory: "AI Inference"
---

# hsdp_ai_inference_compute_environment

Manages HSDP AI Inference Compute Environments

## Example usage

```hcl
data "hsdp_config" "inference" {
  service = "inference"
}

data "hsdp_ai_inference_service_instance" "inference" {
  base_url        = data.hsdp_config.inference.url
  organization_id = var.inference_tenant_org_id
}

resource "hsdp_ai_inference_compute_environment" "compute" {
  endpoint = data.hsdp_ai_inference_service_instance.inference.endpoint
  
  name  = "python3.8_keras_gpu"
  image = "arn:aws:ecr:us-west-2:012345678910:repository/test"
}
```

The following arguments are supported:

* `endpoint` - (Required) The AI Inference instance endpoint
* `name` - (Required) The name of Compute Environment
* `image` - (Required) The image to use for the Compute Environment

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the Compute Environment
* `reference` - The reference of this Compute Environment
* `is_factory` - Weather this Compute Environment is factory one
* `created` - The date this Compute Environment was created
* `created_by` - Who created the environment

## Import

An existing Compute Environment can be imported using `terraform import hsdp_ai_inference_compute_environment`, e.g.

```bash
terraform import hsdp_ai_inference_compute_environment.env a-guid
