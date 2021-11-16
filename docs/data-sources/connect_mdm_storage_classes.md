---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_storage_classes

Retrieve details of Storage classes

## Example Usage

```hcl
data "hsdp_connect_mdm_storage_classes" "all" {
}
```

```hcl
output "storage_class_names" {
   value = data.hsdp_connect_mdm_storage_classes.all.names
}
```

## Attributes Reference

The following attributes are exported:

* `ids` - The storage class IDs
* `names` - The Storage class names
* `descriptions` - The Storage class descriptions
