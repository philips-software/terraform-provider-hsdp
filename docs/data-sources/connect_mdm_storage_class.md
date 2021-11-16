---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_storage_class

Retrieve details of a Storage class

## Example Usage

```hcl
data "hsdp_connect_mdm_storage_class" "postgres" {
  name = "postgres"
}
```

```hcl
output "postgres_storage_class_id" {
   value = data.hsdp_connect_mdm_storage_class.postgres.id
}
```

## Attributes Reference

The following attributes are exported:

* `ids` - The storage class IDs
* `names` - The Storage class names
* `descriptions` - The Storage class descriptions
