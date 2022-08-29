---
subcategory: "DICOM Store"
---

# hsdp_dicom_repository

This resource manages a DICOM repository

## Example Usage

```hcl
resource "hsdp_dicom_repository" "repo1" {
  config_url = hsdp_dicom_store_config.dicom.config_url
  organization_id = hsdp_iam_org.root_org.id
  object_store_id = hsdp_dicom_object_store.store1.id
  
  notification {
    organization_id = hsdp_iam_org.tenant1.id   
  }
}
```

## Argument reference

* `config_url` - (Required) The base config URL of the DICOM Store instance
* `organization_id` - (Required) The organization ID
* `object_store_id` - (Required) the Object store ID
* `store_as_composite` - (Optional) Configure this repository as store as composite.
* `repository_organization_id` - (Optional) The organization ID attached to this repository.
  When not specified, the root organization is used.
* `notification` (Block, Optional)
  * `enabled` - (Required) Enable notifications or not. Default: `true`
  * `organization_id` - (Required) the tenant IAM Organization ID
