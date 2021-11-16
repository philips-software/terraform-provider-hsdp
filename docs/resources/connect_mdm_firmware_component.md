---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_firmware_component

Create and manage MDM FirmwareComponent resources

## Example Usage

```hcl
resource "hsdp_connect_mdm_firmware_component" "first" {
  name                   = "tf-firmware-component"
  description            = "Terraform managed firmware component"
  
  device_type_id = hsdp_connect_mdm_device_type.first.id
  main_component = true
}
```

## Attributes Reference

The following attributes are exported:

* `name` - (Required) The name of the device group
* `description` - (Optional) A short description of the device group
* `device_type_id` - (Required) Reference to the DeviceType
* `main_component` - (Required) Signals if this is a main component (default: `true`)

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `Group/${GUID}`)
* `guid` - The GUID of the service action
