---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_device_type

Create and manage MDM DeviceType resources

## Example Usage

```hcl
resource "hsdp_connect_mdm_device_type" "some_device_type" {
  name                   = "some-device-type"
  description            = "WEARABLE0001"
  commercial_type_number = "WATCH1"
  
  device_group_id = hsdp_connect_mdm_device_group.some_group.id
  
  default_iam_group_id = data.hsdp_iam_group.wearable_group.id
  
  custom_type_attributes = {
    position = "wrist"
    region   = "eu"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the device type
* `description` - (Optional) A short description of the device type
* `device_group_id` - (Required) Reference to the Device Group this type falls under
* `commercial_type_number` - (Required) Commercial Type Number
* `default_iam_group_id` - (Optional) The IAM Group from which this group will inherit roles from
* `custom_type_attributes` - (Optional) Type attributes for all devices under this type.

~> The `name` maps to an AWS IoT thing type so this should be globally unique and not used (or re-used) across deployments

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `Group/${GUID}`)
* `guid` - The GUID of the service action
