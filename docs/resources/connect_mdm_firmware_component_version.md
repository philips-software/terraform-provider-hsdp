---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_firmware_component_version

Create and manage MDM FirmwareComponentVersion resources

## Example Usage

```hcl
resource "hsdp_connect_mdm_firmware_component_version" "one_dot_oh" {
  version                = "1.0.0"
  effective_date         = "2021-10-28"
  description            = "Terraform managed firmware component version"
  
  firmware_component_id = hsdp_connect_mdm_firmware_component.main.id
  component_required    = true
  blob_url              = "/release/1.0.0/firmware.bin"  
  size                  = 512000
  
  fingerprint {
    algorithm = "SHA-256"
    hash = "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
  }
  
  encryption_info {
    encrypted = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `version` - (Required) The version of the Firmware Component image
* `effective_date` - (Required) The effective date of this firmware (Format: `yyyy-mm-dd`)
* `description` - (Optional) A short description of the resource
* `firmware_component_id` - (Required) Reference to Firmware Component resource
* `blob_url` - (Optional) The path of the image on Blob storage
* `size` - (Optional, int) The size of the image
* `component_required` - (Optional, bool) Is the component required (default: `false`)
* `custom_resource` - (Optional, string) JSON string describing your custom resource
* `deprecated_date` - (Optional, date) Deprecated date of this firmware
* `fingerprint` - (Optional) Fingerprint information
  * `algorithm` - (Required) The algorithm used to calculate the fingerprint
  * `hash` - (Required) The fingerprint value
* `encryption_info` - (Optional) Specify encrypted related info
  * `encrypted` - (Required, bool) If the component is encrypted
  * `algorithm` - (Optional) The encryption algorithm that is used
  * `decryption_key` - (Optional) The decryption key

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `FirmwareComponentVersion/${GUID}`)
* `guid` - The GUID of the service action
