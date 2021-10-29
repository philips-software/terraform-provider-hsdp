---
subcategory: "Container Host"
---

# hsdp_container_host_instances

Retrieve a list of container hosts instances

## Example Usage

```hcl
data "hsdp_container_host_instances" "all" {
} 

output "all_container_hosts" {
  value = data.hsdp_container_host_instances.all.ids
}
```

## Attributes Reference

The following attributes are exported:

* `ids` -  The list of container host IDs
* `names` - The list of container host names. This matches up with the `ids` list index.
* `types` - The list of container host instance types. This matches up with the `ids` list index.
* `owners` - The list of container host owners. This matches up with the `ids` list index.
* `private_ips` - The list of container host private IPs. This matches up with the `ids` list index.
