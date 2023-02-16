---
subcategory: "Identity and Access Management (IAM)"
page_title: "HSDP: hsdp_iam_group_membership"
description: |-
  Manages HSDP IAM Group membership resources
---

# hsdp_iam_group_membership

Provides a resource for managing IAM Group membership of users and services.
This resource is useful when the IAM Group is defined or managed elsewhere and
you want to manage membership of a subset of users or services.

~> Use this resource sparingly and carefully, it's easy to create perma-diffs if Terraform declarations conflict. If the IAM group is managed in Terraform make sure `drift_detection` is disabled in its declaration. Also note that the calling identity needs `GROUP.WRITE` access in the groups' managing organization.

## Example Usage

The following example adds users to a group obtained via a data source

```hcl
resource "hsdp_iam_group_membership" "remote_developers" {
  iam_group_id = data.hsdp_iam_group.developers.id
  users        = [hsdp_iam_user.developer1.id, hsdp_iam_user.developer1.id]
}
```

## Argument Reference

The following arguments are supported:

* `iam_group_id` - (Required) The ID of the IAM Group to add users or services to
* `users` - (Optional) The list of user IDs to include in this group. The provider only manages this list of users.
* `services` - (Optional) The list of service identity IDs to include in this group.

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the group membership resource
