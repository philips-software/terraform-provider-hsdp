# hsdp_inference_instance

Retrieve details of an existing HSDP AI Inference instance.

## Example Usage

```hcl
data "hsdp_config" "inference" {
  service = "inference"
}

data "hsdp_inference_instance" "inference" {
  base_url        = data.hsdp_config.inference.url
  organization_id = var.inference_tenant_org_id
}
```

## Argument Reference

The following arguments are supported:

* `base_url` - (Required) the base URL of the Inference deployment. This can be auto-discovered and/or provided by HSDP.
* `organization_id` - (Required) the Tenant IAM organization.

## Attributes Reference

The following attributes are exported:

* `endpoint` - The Inference endpoint URL
