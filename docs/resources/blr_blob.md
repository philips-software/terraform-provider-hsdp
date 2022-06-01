---
subcategory: "Blob Repository (BLR)"
---

# hsdp_blr_blob

Provides a resource for managing [Blob](https://www.hsdp.io/documentation/blob-repository) metadata objects

## Example Usage

The following example creates a Blob metadata object

```hcl
resource "hsdp_blr_blob" "firmware" {
  data_type_name = data.hsdp_mdm_data_type.firmware.name

  blob_path    = "/fw/1.0"
  blob_name    = "firmware-1.0"

  tags = {
    region = "eu-west"
  }
  
  attachment {
    content_type = "application/firmware"  
    data = filebase64(var.firmware_file)
  }
  
  policy {
    statement {
      principal {
        hsdp = ["string"]
      }
      effect = "Allow"
      action = ["GET", "PUT"] 
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `data_type_name` - (Required) The data type name. 
* `blob_path` - (Required) The blob path to use
* `blob_name` - (Required) The blob name to use
* `virtual_path` - (Required) The virtual path to use
* `virtual_name` - (Required) The virtual name to use
* `tags` - (Optional, Hash) A set of tags to associate to this Blob
* `attachment` - (Optional, block) Use this block to define an attachment
  * `content_type` - (Required) The content type of the attachment
  * `data` - (Required) The base64 encoded data
* `policy` - (Optional) Block describing access policy
  * `statement` - (Optional) The policy statement
    * `principal` - (Required) The principal block
      * `hsdp` - (Required, list) The list of `hsdp` principal resource list
    * `effect` - (Required) The Effect element is required and specifies whether
      the statement results in an allow or an explicit deny. The only valid value for Effect is `Allow`.
      Deny is not supported in this API yet.
    * `action` - (Required, list) Specifies the list of actions e.g. s3:GetObject


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
