---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_region

Retrieve details of a region

## Example Usage

```hcl
data "hsdp_connect_mdm_region" "us_east" {
  name = "us-east-1"
}
```

```hcl
output "us_east_region_id" {
   value = data.hsdp_connect_mdm_region.us_east.id
}
```

## Argument Reference

The following arguments are supports:

* `name` - (Required) The name of the region to lookup

## Attributes Reference

The following attributes are exported:

* `id` - The region ID
* `description` - The region description
* `category` - The category of the region
* `hsdp_enabled` - If the regions is enabled
