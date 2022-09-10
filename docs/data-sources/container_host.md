---
subcategory: "Container Host"
---

# hsdp_container_host

Retrieve information from a named Container Host instance

## Example Usage

```hcl
data "hsdp_container_host" "server" {
  name = "my-server.dev"
} 

output "my_server_private_ip" {
  value = data.hsdp_container_host.server.private_ip
}
```

## Argument Reference

* `name` - (Required) The name of the Container Host instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` -  The list of container host IDs
* `type` - The Container Host type
* `owner` - The owner of this Container Host
* `private_ip` - The private IP / address
* `public_ip` - The public IP / address
* `security_groups` - The assigned security groups
* `ldap_groups` - The assigned LDAP groups
* `block_devices` - The provisioned bock devices
* `state` - The state of Container Host instanced
* `subnet` - The subnet where this Container Host is in
* `vpc` - The VPC this Container Host sits in
* `zone` - The network Zone of this Container Host
* `tags` - The tags associated with this Container Host
* `protection` - When set to true delete protection is enabled
