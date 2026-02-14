---
subcategory: "Identity and Access Management (IAM)"
---

# hsdp_iam_application

Retrieve details of an existing application

## Example Usage

The following example looks up an application by ID:

```hcl
data "hsdp_iam_application" "my_app" {
   application_id = var.my_app_id
}
```

The following example looks up an application by name and proposition:

```hcl
data "hsdp_iam_application" "my_app" {
   name = "MYAPP"
   proposition_id = data.hsdp_iam_proposition.my_prop.id
}
```

```hcl
output "my_app_id" {
   value = data.hsdp_iam_application.my_app.id
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Optional) The UUID of the application to look up
* `name` - (Optional) The name of the application to look up
* `proposition_id` - (Optional) the UUID of the proposition the application belongs to

~> When `application_id` is not provided, both `name` and `proposition_id` are required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the application
* `description` - The description of the application
* `global_reference_id` - The global reference ID of the application
