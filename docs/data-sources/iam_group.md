---
subcategory: "Identity and Access Management (IAM)"
---

# hsdp_iam_group

Retrieve details of an existing group

## Example Usage

```hcl
data "hsdp_iam_group" "some_group" {
   managing_organization_id = var.org_id
   name = "GROUP NAME"
}
```

```hcl
output "group_guid" {
   value = data.hsdp_iam_group.id
}
```

## Argument Reference

The following arguments are supported:

* `managing_organization_id` - (Required) the UUID of the managing organization of the group to lookup0
* `name` - (Required) The name of the group to look up

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The group GUID
* `description` - The description of the group
