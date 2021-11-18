---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_firmware_distribution_request

Create and manage MDM FirmwareDistributionRequest resources

## Example Usage

```hcl
resource "hsdp_connect_mdm_firmware_distribution_request" "distro" {
  firmware_version  = "1.0.0"
  description       = "Terraform managed firmware distribution request"
  
  status = "ACTIVE"
  
  distribution_target_device_groups_ids = [
    hsdp_connect_mdm_device_group.first.id,
    hsdp_connect_mdm_device_group.second.id
  ]
  
  firmware_component_version_ids = [
    hsdp_connect_mdm_firmware_component_version.one_dot_oh.id
  ]
  
  orchestration_mode    = "continuous"
  user_consent_required = false
}
```

## Argument Reference

The following arguments are supported:

* `firmware_version` - (Required) The version of the Firmware Component image
* `status` - (Required) The status of the request [`ACTIVE` | `CANCELED`]
* `description` - (Optional) A short description of the resource
* `distribution_target_device_groups_ids` - (Required, list(string)) Reference to Firmware Component resource
* `firmware_component_version_ids` - (Requires, list(string)) The path of the image on Blob storage
* `orchestration_mode` - (Required) What mode of orchestration to use [`none` | `continuous` | `snapshot`]
* `user_consent_required` - (Optional, bool) Is user consent needed for this update (default: `false`)

~> The status field can only be changed to `CANCELED`. This resource is also deprecated, so use it cautiously

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `FirmwareDistributionRequest/${GUID}`)
* `guid` - The GUID of the service action
  
