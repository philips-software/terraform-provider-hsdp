---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_application

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

* `id` - The ID of the application
* `guid` - The raw GUID of the application (MDM reference)
* `description` - The description of the application
* `global_reference_id` - The global reference ID of the application
* `application_guid_system` - The external system associated with resource (this would point to an IAM deployment)
* `application_guid_value` - The external value associated with this resource (this would be an underlying IAM application ID)
* `default_group_guid_system` - The default group guid system set for this resource
* `default_group_guid_value` - The default group guid value set for this resource (IAM group ID)
