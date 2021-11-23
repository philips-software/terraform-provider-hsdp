---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_bucket

Create and manage MDM Bucket resources

## Example Usage

```hcl
resource "hsdp_connect_mdm_bucket" "store" {
  name        = "bucket-store-1"
  description = "Terraform provisioned bucket"
 
  proposition_id    = data.hsdp_connect_mdm_propososition.prop.id
  default_region_id = data.hsdp_connect_mdm_region.us_east.id
  
  versioning_enabled = true
  logging_enabled    = true
  auditing_enabled   = true
  
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
* `description` - (Optional) A short description of the device group
* `proposition_id` - (Required) The proposition ID where this bucket falls under
* `default_region_id` - (Required) The Region this bucket should be provisioned in
* `versioning_enabled` - (Required) Set versioning
* `logging_enabled` - (Required) Set logging
* `auditing_enabled` - (Required) Set auditing
* `enable_cdn` - (Optional) Enable CDN or not
* `cache_control_age` - (Optional) Cache control age settings (Max: `1800`, Min: `300`, Default: `300`)
* `cors_configuration` - (Optional)
  * `allowed_origins` - (Required, list(string)) List of allowed origins
  * `allowed_methods` - (Required, list(string)) Allowed methods: [`GET`, `PUT`, `POST`, `DELETE`, `HEAD`]
  * `max_age_seconds` - (Optional) Max age in seconds
  * `expose_headers` - (Optional, list(string)) List of headers to expose
  
## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `Bucket/${GUID}`)
* `guid` - The GUID of the bucket
