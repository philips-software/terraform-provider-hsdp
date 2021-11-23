---
subcategory: "Identity and Access Management (IAM)"
---

# hsdp_iam_group

Provides a resource for managing HSDP IAM groups

## Example Usage

The following example creates a group

```hcl
resource "hsdp_iam_group" "tdr_users" {
  managing_organization = hsdp_iam_org.devorg.id
  name                  = "TDR Users"
  description           = "Group for TDR Users with Contract and Dataitem roles"
  roles                 = [hsdp_iam_role.TDRALL.id]
  users                 = [hsdp_iam_user.admin.id, hsdp_iam_user.developer.id]
  services              = [hsdp_iam_service.test.id]
}
```

This assumes a role definition exists example:

```hcl
resource "hsdp_iam_role" "TDRALL" {
  name        = "TDRALL"
  description = "Role for TDR users with ALL access"

  permissions = [
    "DATAITEM.CREATEONBEHALF",
    "DATAITEM.READ",
    "DATAITEM.DELETEONBEHALF",
    "DATAITEM.DELETE",
    "CONTRACT.CREATE",
    "DATAITEM.READONBEHALF",
    "CONTRACT.READ",
    "DATAITEM.CREATE",
  ]

  managing_organization = hsdp_iam_org.devorg.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the group
* `description` - (Required) The description of the group
* `roles` - (Required) The list of role IDS to assign to this group
* `managing_organization` - (Required) The managing organization ID
* `users` - (Optional) The list of user IDs to include in this group. The provider only manages this list of users. Existing users added by others means to the group by the provider. It is not practical to manage hundreds or thousands of users this way of course.
* `services` - (Optional) The list of service identity IDs to include in this group. See `hsdp_iam_service`
* `drift_detection` - (Optional, bool) While most resources do automatic drift detection, we are opting to make this
  opt-in feature for IAM Groups due to insufficient IAM API capabilities to perform this operation efficiently.
  A future version might change this to be always-on. When enabled, the provider will perform additional API calls
  to determine if any changes were made outside of Terraform to user and service assignments of this Group. Default: `false`

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the group

## Import

An existing group can be imported using `terraform import hsdp_iam_group`, e.g.

```shell
terraform import hsdp_iam_group.mygroup a-guid
```
