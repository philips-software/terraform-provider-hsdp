# hsdp_iam_permission
Provides a resource for managing HSDP IAM permissions

## Example Usage

The following example creates a new permission

```hcl
resource "hsdp_iam_permission" "NEWAPP_ACCESS" {
  name = "NEWAPP.ACCESS"
  description = "permission to access newapp data"
  category = "TEST"
  type = "GLOBAL"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the permission
* `description` - (Required) Description of the permission
* `category` - (Required) Category of the permission
* `type` - (Required) The type. Supported: GLOBAL

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the permission

## Import

An existing proposition can be imported using `terraform import hsdp_iam_permission`, e.g.

```shell
terraform import hsdp_iam_permission.mypermission a-guid
```

