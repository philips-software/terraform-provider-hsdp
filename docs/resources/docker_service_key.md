# hsdp_docker_service_key

Manages HSDP Docker Registry service keys

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
* `username` - (Computed) The service id
* `password` - (Computed, Sensitive) The active private of the service
* `created_At` - (Computed) The timestamp this key was created
