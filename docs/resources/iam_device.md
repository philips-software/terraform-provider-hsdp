---
subcategory: "Identity and Access Management (IAM)"
---

# hsdp_iam_device

Provides a resource for managing HSDP IAM devices.

## Example Usage

The following example creates a device

```hcl
resource "random_password" "test_device" {
}

resource "random_uuid" "test_device" {
}

// Create a test device
resource "hsdp_iam_device" "test_device" {
  login_id = "test_device_a"
  password = random_password.test_device.result
  
  external_identifier {
    type {
      code = "ID"
      text = "Device Identifier"
    }
    system = "http://www.philips.co.id/c-m-ho/cooking/airfryer"
    value = "001"
  }
  
  
  type = "ActivityMonitor"
  text = "This is a test device managed by Terraform"
  
  organization_id = var.org_id
  application_id  = var.app_id
  
  for_test  = true
  is_active = true
  
  global_reference_id = random_uuid.test_device.result
}
```

## Argument Reference

The following arguments are supported:

* `login_id` - (Required) The login id of the device
* `password` - (Required) The password of the device
* `external_identifier` - (Required) Block describing external ID of this device
  * `type` - (Required) - Block describing the type
    * `code` - (Required) The code of the ID
    * `text` - (Optional) Text describing the code type
  * `system` - (Required) The system where the identifier is defined in
  * `value` - (Required) The value of the identifier
* `organization_id` - (Required) the organization ID (GUID) this device should be attached to
* `application_id` - (Required) the application ID (GUID) this device should be attached to
* `for_test` - (Optional) Boolean. When set to true this device is marked as a test device
* `is_active` - (Optional) Boolean. Controls if this device is active or not

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the device
* `registration_date` - (Generated) The date the device was registered

## Import

Existing devices can be imported, however they will be missing their password rendering them pretty much useless.
Therefore, we recommend creating them using the provider.
