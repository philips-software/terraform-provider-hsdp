---
subcategory: "Docker Registry"
---

# hsdp_docker_namespace

Manage HSDP Docker registry namespaces

~> This resource only works when `HSDP_UAA_USERNAME` and `HSDP_UAA_PASSWORD` values matching provider arguments are set

## Example usage

```hcl
resource "hsdp_docker_namespace" "project1" {
  name = "project1"
  description = "project1 namespace"
  full_description = ""
}
```

## Argument reference

* `name` - (Required) The name of the namespace to look up

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the namespace
