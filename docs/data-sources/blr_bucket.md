---
subcategory: "Blob Repository (BLR)"
---

# hsdp_blr_bucket

Retrieve details on a Blob Repository Bucket resource

## Example Usage

```hcl
data "hsdp_blr_bucket" "store" {
  name        = "bucket-store-1"
}

output "bucket_id" {
  value = data.hsdp_connect_mdm_bucket.store.id
}
```

## Argument Reference

The following arguments are available:

* `name` - (Required) The name of the bucket to look up

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `Bucket/${GUID}`)
* `guid` - The GUID of the bucket
* `cdn_enabled` - CDN enabled or not
* `cache_control_age` - Cache control age settings
* `cors_config_json` - The Bucket CORS configuration in JSON
