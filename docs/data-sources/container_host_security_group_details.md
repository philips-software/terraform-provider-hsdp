---
subcategory: "Container Host"
---

# hsdp_container_host_security_group_details

Provides details of a specific security group rules.

> This data source is only available when the `cartel_*` keys are set in the provider config

## Example Usage

```hcl
data "hsdp_container_host_security_group_details" "http_from_cf" {
  name = "http-from-cloud-foundry"
}

output "port_ranges" {
   value = data.hsdp_container_host_security_group_details.http_from_cf.port_ranges
}
```

## Attributes Reference

The following attributes are exported:

* `port_ranges` - The port ranges associated to the rule
* `protocols` - The protocol of the rule
* `sources` - The source address of the rule
