# hsdp_container_host
Provides HSDP Container Host instances

> This resource is only available when the `cartel_*` keys are set in the provider config

## Example Usage

The following example provisions three (3) new container host instances

```hcl
resource "hsdp_container_host" "zahadoom" {
  count         = 3
  name          = "zahadoom-${count.index}.dev"
  volumes       = 1
  volume_size   = 50
  instance_type = "t2.medium"

  user_groups     = var.user_groups
  security_groups = ["analytics", "tcp-8080"]

  connection {
    bastion_host = var.bastion_host
    host         = self.private_ip
    user         = var.user
    private_key  = var.private_key
    script_path  = "/home/${var.user}/bootstrap.bash"
  }

  tags {
    "created_by" = "terraform"
  }

  provisioner "remote-exec" {
    inline = [
      "ifconfig",
      "docker volume create fluent-bit",
      "docker run -d -p 24224:24224 -v fluent-bit:/fluent-bit/etc philipssoftware/fluent-bit-out-hsdp:1.4.4"
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The container host name. Must be unique.
* `instance_type` - (Optional) The EC2 instance type to use. Default `m5.large`
* `instance_role` - (Optional) The role to use. Default `container-host` (other values: `vanilla`, `base`)
* `volume_type` - (Optional) The EBS volume type.
* `iops` - (Optional) Number of IOPs to provision.
* `protect` - (Optional) Boolean when set will enable protection for container host.
* `encrypt_volumes` - (Optional) When set encrypts volumes. Default is `true`
* `volumes` - (Optional) Number of additional volumes to attach. Default `0`
* `volume_size` - (Optional) Volume size in GB.
* `security_groups` - (Optional) list(string) of Security groups to attach. Default `[]`
* `user_groups` - (Optional) list(string) of User groups to attach. Default `[]`
* `tags` - (Optional) Map of tags to assign to the instances

## Attributes Reference

The following attributes are exported:

* `id` - The instance ID
* `private_ip` - The private IP address of the instance
* `role` - The role of the instance.
* `subnet` - The subnet the instance was provisioned in.
* `vpc` - The VPC the instance was provisioned in.
* `zone` - The Zone the instance was provisioned in.
* `launch_time` - Timestamp when the instance was launched.
* `block_devices` - The list of block devices attached to the instance.

## Import

Importing existing instances is supported but not recommended.
