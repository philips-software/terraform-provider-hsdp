---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_service_actions

Retrieve service actions

## Example Usage

```hcl
data "hsdp_connect_mdm_service_actions" "all" {
}
```

```hcl
output "service_action_names" {
   value = data.hsdp_connect_mdm_service_actions.all.names
}

output "standard_service_ids" {
  value = data.hsdp_connect_mdm_service_actions.all.standard_service_ids
}
```

## Attributes Reference

The following attributes are exported:

* `ids` - The GUID of the proposition
* `names` - The names of the service actions
* `standard_service_ids` - the GUIDs of the standard services associated with the service actions
