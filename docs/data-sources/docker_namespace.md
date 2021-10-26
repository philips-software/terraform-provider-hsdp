# hsdp_docker_namespace

Retrieve details on a HSDP Docker Registry namespace

~> This resource only works when `HSDP_UAA_USERNAME` and `HSDP_UAA_PASSWORD` values matching provider arguments are set

## Example usage

```hcl
data "hsdp_docker_namespace" "project1" {
  name = "project1"
}

output "repositories" {
  value = data.hsdp_docker_namespace.project1.repositories
}

```

## Argument reference

* `name` - (Required) The name of the namespace to look up

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the namespace
* `repositories` - (list(string)) The list of repositories in this namespace
