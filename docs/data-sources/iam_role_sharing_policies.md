---
subcategory: "Identity and Access Management (IAM)"
---

# hsdp_iam_role_sharing_policies

Lists defined IAM Role Sharing Policies

## Example Usage

```hcl
data "hsdp_iam_role_sharing_policies" "all" {
  role_id = var.role_id
}
```

```hcl
output "policy_ids" {
   value = data.hsdp_iam_role_sharing_policies.all.ids
}
```

## Argument Reference

The following arguments are supported:

* `role_id` - (Required) The role ID to search for policies
* `sharing_policy` - (Optional) Only list results with this policy sharing
* `target_organization_id` - (Optional) Only list policies targetting the given organization

## Attributes Reference

The following attributes are exported:

* `ids` - The IDS of the sharing policies
* `role_names` - The role names
* `sharing_policies` - The sharing policies
* `target_organization_ids` - The target organization IDs
* `source_organization_ids` - The source organization IDs
* `purposes` - The purposes
