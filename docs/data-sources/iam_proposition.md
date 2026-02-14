---
subcategory: "Identity and Access Management (IAM)"
---

# hsdp_iam_proposition

Retrieve details of an existing proposition

## Example Usage

The following example looks up a proposition by ID:

```hcl
data "hsdp_iam_proposition" "my_prop" {
   proposition_id = var.my_prop_id
}
```

The following example looks up a proposition by name and organization:

```hcl
data "hsdp_iam_proposition" "my_prop" {
   name = "MYPROPOSITION"
   organization_id = var.my_org_id
}
```

```hcl
output "my_prop_display_name" {
   value = data.hsdp_iam_proposition.my_prop.display_name
}
```

## Argument Reference

The following arguments are supported:

* `proposition_id` - (Optional) The UUID of the proposition to look up
* `name` - (Optional) The name of the proposition to look up
* `organization_id` - (Optional) the UUID of the organization the proposition belongs to

~> When `proposition_id` is not provided, both `name` and `organization_id` are required.

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the proposition
* `description` - The description of the proposition
* `global_reference_id` - The global reference ID of the proposition
