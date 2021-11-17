---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_device_group

Create and manage MDM DeviceGroup resources

## Example Usage

```hcl
resource "hsdp_connect_mdm_device_group" "some_group" {
  name        = "some-device-group"
  description = "A device group"
  
  application_id = data.hsdp_connect_mdm_application.app.id
  
  default_iam_group_id = data.hsdp_iam_group.device_group.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the device group
* `description` - (Optional) A short description of the device group
* `application_id` - (Required) Reference to the Application this group falls under
* `default_iam_group_id` - (Optional) The IAM Group from which this group will inherit roles from

~> The `name` maps to an AWS IoT thing group so this should be globally unique and not used (or re-used) across deployments

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `Group/${GUID}`)
* `guid` - The GUID of the service action
