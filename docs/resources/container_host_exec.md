# hsdp_container_host_exec

Copies content and executes commands on Container Host instances

> This resource is only available when the `cartel_*` keys are set in the provider config

## Example Usage

The following example uses the internal provisioning support for bootstrapping an instance

```hcl
resource "hsdp_container_host_exec" "init" {
  host = hsdp_container_host.mybox.private_ip
  user = var.user
  # an SSH agent should be running with the SSH private key loaded
   
  file {
    content = "echo Hello world"
    destination = "/tmp/hello.sh"
    permissions = "0755"
  }
  
  commands = [
    "/bin/sh -C /tmp/hello.sh"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `user` - (Required) The username to use for provision activities using SSH
* `private_key` - (Optional) The SSH private key to use for provision activities. When not provided an ssh-agent should be available.
* `file` - (Optional) Block specifying content to be written to the container host after creation
* `commands` - (Required, list(string)) List of commands to execute after creation of container host
* `bastion_host` - (Optional) The bastion host to use.  When not set, this will be deduced from the container host location
* `triggers` - (Optional, list(string)) An list of strings which when changes will trigger recreation of the resource triggering
   all create files and commands executions.

Each `file` block can contain the following fields. Use either `content` or `source`:

* `source` - (Optional, file path) Content of the file. Conflicts with `content`
* `content` - (Optional, string) Content of the file. Conflicts with `source`
* `destination` - (Required, string) Remote filename to store the content in
* `permissions` - (Optional, string) The file permissions. Default permissions are "0644"
* `owner` - (Optional, string) The file owner. Default owner the SSH user
* `group` - (Optional, string) The file group. Default group is the SSH user's group

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID
