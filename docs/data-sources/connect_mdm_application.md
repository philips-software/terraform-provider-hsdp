---
subcategory: "Master Data Management (MDM)"
---

# hsdp_iam_application

Retrieve details of an existing application

## Example Usage

```hcl
data "hsdp_connect_mdm_application" "app" {
   name = "MYAPP"
   proposition_id = data.hsdp_iam_proposition.app.id
}
```

```hcl
output "my_app_id" {
   value = data.hsdp_connect_mdm_application.app.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the application to look up
* `proposition_id` - (Required) the UUID of the proposition the application belongs to

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the application
* `description` - The description of the application
* `global_reference_id` - The global reference ID of the application
