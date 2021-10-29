---
subcategory: "Docker Registry"
---

# hsdp_docker_namespaces

Retrieves information on available HSDP Docker Registry namespaces

~> This resource only works when `HSDP_UAA_USERNAME` and `HSDP_UAA_PASSWORD` values matching provider arguments are set

## Example usage

```hcl
data "hsdp_docker_namespaces" "namespaces" {
}

output "namespaces" {
  value = data.hsdp_docker_namespaces.namespaces.names
}
```

## Attribute reference

The following attributes are available:

* `names` - The list of available namespaces
* `num_repos` - The number of repositories. Index matches the `names` list
