---
subcategory: "Docker Registry"
---

# hsdp_docker_service_key

Manages HSDP Docker Registry service keys

~> This resource only works when `HSDP_UAA_USERNAME` and `HSDP_UAA_PASSWORD` values matching provider arguments are set.

## Example usage

```hcl
resource "hsdp_docker_registry_key" "cicd" {
  description = "Terraform managed key"
}

output "docker_username" {
  value = hsdp_docker_registry_key.cicd.username
}

output "docker_password" {
  value = hsdp_docker_registry_key.cicd.password
}
```

## Argument reference

The following arguments are supported:

* `description` - (Required) The description of the service key

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the service key
* `username` - (Computed) The service key username
* `password` - (Computed, Sensitive) The service key password
* `created_at` - (Computed) The timestamp the key was created
