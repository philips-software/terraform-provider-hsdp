---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_resource_limits

Retrieve resource limits configured in Connect MDM

## Example Usage

```hcl
data "hsdp_connect_mdm_resource_limits" "all" {
}
```

```hcl
output "limited_resources" {
   value = data.hsdp_connect_mdm_resource_limits.all.resources
}

output "limited_defaults" {
  value = data.hsdp_connect_mdm_resource_limits.all.defaults
}
```

## Attributes Reference

The following attributes are exported:

* `resources` - The region IDs
* `defaults` - the region names
* `overrides` - The region descriptions
