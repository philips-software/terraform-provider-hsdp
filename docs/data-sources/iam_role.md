---
subcategory: "IAM"
---

# hsdp_iam_role

Retrieve details of an existing role

## Example Usage

```hcl
data "hsdp_iam_role" "some_role" {
   managing_organization_id = var.org_id
   name = "GROUP NAME"
}
```

```hcl
output "role_guid" {
   value = data.hsdp_iam_role.id
}
```

## Argument Reference

The following arguments are supported:

* `managing_organization_id` - (Required) the UUID of the managing organization of the role
* `name` - (Required) The name of the role to look up

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The group GUID
* `description` - The description of the group
