---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_data_adapters

Retrieve details of DataAdapters

## Example Usage

```hcl
data "hsdp_connect_mdm_data_adapters" "all" {
}
```

```hcl
output "data_adapters_names" {
   value = data.hsdp_connect_data_adapters.all.names
}
```

## Attributes Reference

The following attributes are exported:

* `ids` - The DataAdapter IDs
* `names` - The DataAdapter names
* `descriptions` - The DataAdapter descriptions
* `service_agent_ids` - The service agent IDs
