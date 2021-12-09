---
subcategory: "DICOM Store"
---

# hsdp_dicom_notification

This resource manages a DICOM notification configurations

## Example Usage

```hcl
resource "hsdp_dicom_notification" "topic" {
  config_url = hsdp_dicom_store_config.dicom.config_url
  organization_id = hsdp_iam_org.root_org.id
  endpoint_url = var.notification_endpoint_url
  
  default_organization_id = hsdp_iam_org.tenant1.id
}
```

## Argument reference

* `config_url` - (Required) The base config URL of the DICOM Store instance
* `organization_id` - (Required) The organization ID
* `endpoint_url` - (Required) The notification endpoint URL. Example: `https://notification-dev.us-east.philips-healthsuite.com`
* `default_organization_id` - (Required) The default organization ID
