---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_bucket

Retrieve details on a MDM Bucket resource

## Example Usage

```hcl
resource "hsdp_connect_mdm_bucket" "store" {
  name        = "bucket-store-1"
  proposition_id    = data.hsdp_connect_mdm_proposition.prop.id
}

output "bucket_id" {
  value = data.hsdp_connect_mdm_bucket.store.id
}
```

## Argument Reference

The following arguments are available:

* `name` - (Required) The name of the bucket to look up
* `proposition_id` - (Required) The proposition ID where this bucket falls under

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `Bucket/${GUID}`)
* `guid` - The GUID of the bucket
* `default_region_id` - The Region this bucket should be provisioned in
* `versioning_enabled` - Versioning enabled
* `logging_enabled` - Logging enabled
* `auditing_enabled` - Auditing enabled
* `cdn_enabled` - CDN enabled or not
* `cache_control_age` - Cache control age settings
* `cors_config_json` - The Bucket CORS configuration in JSON
