Provides details of a given HSDP IAM user.

> This data source is only available when the `cartel_*` keys are set in the provider config

## Example Usage

```hcl
data "hsdp_container_host_subnet_types" "subnets" {
}

output "subnet_names" {
   value = data.hsdp_container_host_subnet_types.subnets.names
}
```

## Attributes Reference

The following attributes are exported:

* `names` - The names of all subnets
* `ids` - Map of ids belonging to names
* `networks` - Map of networks belonging to names
