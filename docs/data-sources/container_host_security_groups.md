---
subcategory: "Container Host"
---

# hsdp_container_host_security_groups

Provides details of the available security groups.

> This data source is only available when the `cartel_*` keys are set in the provider config

## Example Usage

```hcl
data "hsdp_container_host_security_groups" "all" {
}

output "security_groups" {
   value = data.hsdp_container_host_security_groups.all.names
}
```

## Attributes Reference

The following attributes are exported:

* `names` - The names of all available security groups
