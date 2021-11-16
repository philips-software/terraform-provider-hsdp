---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_blob_data_contract

Create and manage MDM BlobDataContract resources

## Example Usage

```hcl
resource "hsdp_connect_mdm_blob_data_contract" "contract" {
  name = "tf-blob-data-contract"

  data_type_id     = hsdp_connect_mdm_data_type.first.id
  bucket_id        = hsdp_connect_mdm_bucket.first.id
  storage_class_id = data.connect_mdm_storage_class.first.id

  root_path_in_bucket = "/"
  
  logging_enabled                  = true
  cross_region_replication_enabled = false
}
```

## Attributes Reference

The following attributes are exported:

* `name` - (Required) The name of the device group
* `data_type_id` - (Required) Reference to the DataType
* `bucket_id` - (Required) Reference to the Bucket
* `storage_class_id` - (Required) Reference to the StorageClass
* `root_path_in_bucket` - (Required) The root path in the bucket
* `logging_enabled` - (Optional) Enable logging (default: `true`)
* `cross_region_replication_enabled` - (Optional) cross region replication active (default: `false`)

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `Group/${GUID}`)
* `guid` - The GUID of the service action
