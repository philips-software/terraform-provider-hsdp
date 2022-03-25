---
subcategory: "Blob Repository (BLR)"
---

# hsdp_blr_blob

Provides a resource for managing [Blob Repository](https://www.hsdp.io/documentation/blob-repository) objects.

## Example Usage

The following example creates a Blob

```hcl
resource "hsdp_blr_blob" "firmware" {
  data_type_id = data.hsdp_mdm_data_type.firmware.id
  blob_path    = "/fw/1.0"
  blob_name    = "firmware-1.0"

  tags = {
    region = "eu-west"
  }
  
  attachment {
    content_type = "application/firmware"  
    data = filebase64(var.firmware_file)
  }
}
```

## Argument Reference

The following arguments are supported:

* `data_type_id` - (Required) The data type ID 
* `blob_path` - (Required) The blob path to use
* `blob_name` - (Required) The blob name to use
* `virtual_path` - (Required) The virtual path to use
* `virtual_name` - (Required) The virtual name to use
* `tags` - (Optional, Hash) A set of tags to associate to this Blob
* `attachment` - (Optional, block) Use this block to define an attachment
  * `content_type` - (Required) The content type of the attachment
  * `data` - (Required) The base64 encoded data
  


## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the Blob repo
* `bucket` - The bucket associated with this Blob
* `creation` - The creation date of this Blob
* `created_by` - Which entity created this Blob

## Import

An existing Blob can be imported using `terraform import hsdp_blr_blob`, e.g.

```bash
terraform import hsdp_blr_blob.blobby a-blob-guid-here
```
