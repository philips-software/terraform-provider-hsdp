# hsdp_container_host

Manage HSDP Container Host instances

> This resource is only available when the `cartel_*` keys are set in the provider config

## Example Usage

The following example provisions and bootstraps a container host instance:

```hcl
resource "hsdp_container_host" "zahadoom" {
  name = "mybox.dev"
  instance_type = "t2.medium"

  user_groups = var.user_groups
  security_groups = ["analytics", var.user]

  bastion_host = var.bastion_host
  user = var.user
  private_key = var.private_key

  tags = {
    created_by = "terraform"
  }

  file {
    content = "This string will be stored remotely"
    destination = "/tmp/stored.txt"
    permissions = "0700"
  }
  
  commands = [
    "cat /tmp/stored.txt",
    "docker volume create fluent-bit",
    "docker run -d -p 24224:24224 -v fluent-bit:/fluent-bit/etc philipssoftware/fluent-bit-out-hsdp:1.4.4"
  ]
}
```

The following example provisions three (3) new container host instances and using Terraform's traditional provisioners

```hcl
resource "hsdp_container_host" "zahadoom" {
  count = 3
  name = "zahadoom-${count.index}.dev"
  volumes = 1
  volume_size = 50
  instance_type = "t2.medium"

  user_groups = var.user_groups
  security_groups = ["analytics", var.user]

  connection {
    bastion_host = var.bastion_host
    host = self.private_ip
    user = var.user
    private_key = var.private_key
    script_path = "/home/${var.user}/bootstrap.bash"
  }

  tags = {
    created_by = "terraform"
    owner = var.user
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
* `image` - (Optional) The OS image to use. Only use this if you have access to additional image types (example: `centos7`). Conflicts with `instance_role` value `container-host`
* `volume_type` - (Optional) The EBS volume type. Default is `gp2`. You can also choose `io1` which is default when you specify `iops` value
* `iops` - (Optional) Number of guaranteed IOPs to provision. Supported value range `1-4000`
* `protect` - (Optional) Boolean when set will enable protection for container host.
* `encrypt_volumes` - (Optional) When set encrypts volumes. Default is `true`
* `volumes` - (Optional) Number of additional volumes to attach. Default `0`, Maximum `6`
* `volume_size` - (Optional) Volume size in GB. Supported value range `1-16000` (16 TB max)
* `security_groups` - (Optional) list(string) of Security groups to attach. Default `[]`, Maximum `4`
* `user_groups` - (Optional) list(string) of User groups to attach. Default `[]`, Maximum `50`
* `subnet` - (Optional) This will cause a new instance to get deployed on a specific subnet. Conflicts with `subnet_type`. You should only use this option if you have very specific requirements that dictate all the instances you are creating need to reside in the same AZ. An example of this would be a cluster of systems that need to reside in the same datacenter.
* `subnet_type` - (Optional) What subnet type to use. Can be `public` or `private`. Default is `private`.
* `tags` - (Optional) Map of tags to assign to the instances
* `user` - (Optional) The username to use for provision activities using SSH
* `private_key` - (Optional) The SSH private key to use for provision activities
* `file` - (Optional) Block specifying content to be written to the container host after creation
* `bastion_host` - (Optional) The bastion host to use.  When not set, this will be deduced from the container host location
* `keep_failed_instances` - (Optional) Keep instances around for post-mortem analysis on failure. Default is `false`.

Each `file` block can contain the following fields. Use either `content` or `source`:

* `source` - (Optional, file path) Content of the file. Conflicts with `content`
* `content` - (Optional, string) Content of the file. Conflicts with `source`
* `destination` - (Required, string) Remote filename to store the content in
* `permissions` - (Optional, string) The file permissions. Default permissions are "0644"
* `owner` - (Optional, string) The file owner. Default owner the SSH user
* `group` - (Optional, string) The file group. Default group is the SSH user's group
* `commands` - (Optional, list(string)) List of commands to execute after creation of container host

-> We recommend using a [hsdp_container_host_exec](https://registry.terraform.io/providers/philips-software/hsdp/latest/docs/resources/container_host_exec) resource to provision files and commands on your instance. This decouples software bootstrapping from the instance provisioning, which can take between 5-15 minutes on its own.

## Attributes Reference

The following attributes are exported:

* `id` - The instance ID
* `private_ip` - The private IP address of the instance
* `public_ip` - The public IP address of the instance if it has one
* `role` - The role of the instance.
* `subnet` - The subnet the instance was provisioned in.
* `vpc` - The VPC the instance was provisioned in.
* `zone` - The Zone the instance was provisioned in.
* `launch_time` - Timestamp when the instance was launched.
* `block_devices` - The list of block devices attached to the instance.

## Import

Importing existing instances is supported but not recommended.
