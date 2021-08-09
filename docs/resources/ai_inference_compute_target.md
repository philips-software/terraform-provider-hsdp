# hsdp_ai_inference_compute_target

Manages HSDP AI Inference Compute Targets

## Example usage

```hcl
data "hsdp_config" "inference" {
  service = "inference"
}

data "hsdp_ai_inference_instance" "inference" {
  base_url        = data.hsdp_config.inference.url
  organization_id = var.inference_tenant_org_id
}

resource "hsdp_ai_inference_compute_target" "target" {
  endpoint = data.hsdp_ai_inference_instance.inference.endpoint
  
  name          = "gpu1"
  description   = "Tesla v100 GPU based environment with 128MB GPU memory"
  instance_type = "ml.p3.16xlarge"
  storage       = 20
}
```

The following arguments are supported:

* `endpoint` - (Required) The AI Inference instance endpoint
* `name` - (Required) The name of Compute Environment
* `instance_type` - (Required) The instance type to use. Available instance types for Inference, see https://aws.amazon.com/sagemaker/pricing/
* `storage` - (Required) Additional storage to allocate (in GB). Default: `1`
* `description` - (Optional) Description of the Compute Target
## Attributes Reference

In addition to all arguments above, the following attributes are exported:
attributes are exported:

* `id` - The GUID of the Compute Target
* `is_factory` - Weather this Compute Environment is factory one
* `created` - The date this Compute Environment was created
* `created_by` - Who created the environment

## Import

An existing Compute Environment can be imported using `terraform import hsdp_ai_inference_compute_target`, e.g.

```bash
terraform import hsdp_ai_inference_compute_target.target a-guid
