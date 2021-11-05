---
subcategory: "Identity and Access Management (IAM)"
---

# hsdp_iam_permission

Retrieves all available permissions

## Example Usage

```hcl
data "hsdp_iam_permissions" "list" {}
```

```hcl
output "all_permissions" {
   value = data.hsdp_iam_permissions.list.permissions
}
```

## Attributes Reference

The following attributes are exported:

* `permissions` - The list of permissions
