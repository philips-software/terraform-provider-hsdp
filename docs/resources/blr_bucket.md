---
subcategory: "Blob Repository (BLR)"
page_title: "HSDP: hsdp_blr_bucket"
description: |-
  Manages HSDP Connect Blob Repository Buckets
---

# hsdp_blr_bucket

Create and manage Blob Repository Buckets

## Example Usage

```hcl
resource "hsdp_blr_bucket" "store" {
  name        = "bucket-store-1"
 
  proposition_id = data.hsdp_connect_mdm_propososition.prop.id
  
  enable_cdn = false
  
  cors_configuration {
    allowed_origins = ["https://foo.hsdp.io"]
    allowed_methods = ["GET"]
    expose_headers  = ["X-Hsdp-Signature"]
  }
}
```

## Argument Reference

The following arguments are available:

* `name` - (Required) The name of the device group
* `proposition_id` - (Required) The proposition ID where this bucket falls under
* `versioning_enabled` - (Required) Set versioning
* `enable_cdn` - (Optional) Enable CDN or not
* `enable_create_or_delete_blob_meta` - (Optional) Enables creation or deletion of Blob meta data
* `enable_hsdp_domain` - (Optional) Enable HSDP domain mapping
* `cache_control_age` - (Optional) Cache control age settings (Max: `1800`, Min: `300`, Default: `1`)
* `cors_configuration` - (Optional)
  * `allowed_origins` - (Required, list(string)) List of allowed origins
  * `allowed_methods` - (Required, list(string)) Allowed methods: [`GET`, `PUT`, `POST`, `DELETE`, `HEAD`]
  * `max_age_seconds` - (Optional) Max age in seconds
  * `expose_headers` - (Optional, list(string)) List of headers to expose

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `Bucket/${GUID}`)
* `guid` - The GUID of the bucket
