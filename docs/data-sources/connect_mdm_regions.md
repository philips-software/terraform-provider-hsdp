---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_regions

Retrieve details of an existing regions

## Example Usage

```hcl
data "hsdp_connect_mdm_regions" "all" {
}
```

```hcl
output "regions" {
   value = data.hsdp_connect_mdm_regions.all.regions
}
```

## Attributes Reference

The following attributes are exported:

* `ids` - The region IDs
* `names` - the region names
* `descriptions` - The region descriptions
* `hsdp_enabled` - Enabled list
