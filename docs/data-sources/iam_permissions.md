---
subcategory: "Identity and Access Management (IAM)"
---

# hsdp_iam_permissions

Retrieves all available permissions

## Example Usage

```hcl
data "hsdp_iam_permissions" "list" {}
```

```hcl
output "all_permissions" {
   value = data.hsdp_iam_permissions.list.names
}

output "permission_descriptions" {
  value = [for i, v in data.hsdp_iam_permissions.all.ids :
    {
      "name" : data.hsdp_iam_permissions.list.names[i],
      "description" : data.hsdp_iam_permissions.list.descriptions[i]
  }]
}
```

## Attributes Reference

The following attributes are exported:

* `ids` - The IDs of the permissions
* `names` - The list of permissions
* `descriptions` - The list of descriptions
* `types` - The types of the permissions
* `categories` - The categories of the permissions
* `permissions` - (Deprecated, use 'names') The name of permissions
