# hsdp_dicom_repository
This resource manages a DICOM repository

# Example Usage
```hcl
resource "hsdp_dicom_repository" "repo1" {
  organization_id = hsdp_iam_org.root_org.id
  object_store_id = hsdp_dicom_object_store.store1.id
}
```

# Argument reference

* `config_url` - (Required) The base config URL of the DICOM Store instance
* `organization_id` - (Required) The organization ID
* `object_store_id` - (Required) the Object store ID
