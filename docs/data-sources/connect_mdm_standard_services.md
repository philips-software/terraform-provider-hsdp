---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_standard_services

Retrieve details of available Standard Services

## Example Usage

```hcl
data "hsdp_connect_mdm_standard_services" "all" {
}
```

```hcl
output "standard_service_ids" {
   value = data.hsdp_connect_mdm_standard_services.all.ids
}
```

## Attributes Reference

The following attributes are exported:

* `ids` - The StandardService IDs
* `names` - The names of the standard services
* `descriptions` - The StandardService descriptions
* `trusted` - If the services are trusted
