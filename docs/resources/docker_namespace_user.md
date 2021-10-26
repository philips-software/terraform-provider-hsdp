# hsdp_docker_namespace_user

Manages user access to namespace. Permissions can be set to pull, push and delete repositories
in the namespace. A user grant also be granted admin permissions

## Example usage

```hcl
resource "hsdp_docker_namespace" "tycho" {
  name = "tycho"
}

resource "hsdp_docker_namespace_user" "gnemath" {
  username     = "camina"
  namespace_id = hsdp_docker_namespace.tycho.id
  
  can_pull   = true
  can_push   = true
  can_delete = true
  is_admin   = false
}
```

Gives user `camina` pull, push and delete access to `tycho` space, but not admin rights.

## Argument reference

The following arguments are supported:

* `username` - (Required) The LDAP / UAA username of the user to configure permissions for
* `namespace_id` - (Required) The namespace ID to configure permissions for
* `can_push` - (Optional) Specifies if the user can push repositories or not. Default: `false`
* `can_pull` - (Optional) Specifies if the user can pull repositories or not. Default: `true`
* `can_delete` - (Optional) Specifies if the user can delete repositories or not. Default: `false`
* `is_admin` - (Optional) Admin permissions on the namespace. Default: `false`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the permission record
* `user_id` - (Computed) The user id
