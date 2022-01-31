---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_application

Create and manage MDM Application resources

~> Currently, deleting Application resources is not supported by the MDM API, so use them sparingly

## Example Usage

```hcl
resource "hsdp_connect_mdm_application" "app" {
  name        = "mobile"
  description = "Terraform managed Application"
  proposition_id = data.hsdp_connect_mdm_proposition.prop.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Application
* `description` - (Optional) A short description of the Application
* `proposition_id` - (Required) The ID of the Proposition this Application should fall under
* `global_reference_id` - (Optional) A global reference ID for this application
* `default_group_guid` - (Optional) The GUID of the IAM Group to assign by default

~> The `proposition_id` only accept MDM Proposition IDs. Using an IAM Proposition ID will not work, even though they might look similar.

~> The `default_group_guid` takes an IAM Group ID i.e. from an `hsdp_iam_group` resource or data source element

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference (format: `Application/${GUID}`)
* `guid` - The GUID of this resource
* `application_guid_system` - The external system associated with resource (this would point to an IAM deployment)
* `application_guid_value` - The external value associated with this resource (this would be an underlying IAM application ID)
