---
subcategory: "DICOM Store"
page_title: "HSDP: hsdp_dicom_notification"
description: |-
  Manages HSDP DICOM Store notifications
---

# hsdp_dicom_notification

This resource manages a DICOM notification configurations

## Example Usage

```hcl
resource "hsdp_dicom_notification" "topic" {
  config_url = hsdp_dicom_store_config.dicom.config_url
  endpoint_url = var.notification_endpoint_url
  
  default_organization_id = hsdp_iam_org.tenant1.id
}
```

## Argument reference

* `config_url` - (Required) The base config URL of the DICOM Store instance
* `endpoint_url` - (Required) The notification endpoint URL. Example: `https://notification-dev.us-east.philips-healthsuite.com`
* `enabled` - (Optional) Enable the notification or not. Default: `true`
* `default_organization_id` - (Optional) The default organization ID
