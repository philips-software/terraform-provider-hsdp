---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_service_actions

Retrieve details of ServiceActions

## Example Usage

```hcl
data "hsdp_connect_mdm_service_actions" "all" {
}
```

## Argument Reference

* `filter` - (Optional) The filter conditions block for selecting ServiceActions

### filter options

* `id` - (Optional) The id (uuid) of the action
* `name` - (Optional) Filter by name
* `organization_guid_value` - (Optional) Filter on organization GUID value
* `standard_service_id` - (Optional) Filter on standard service ID

## Attributes Reference

The following attributes are exported:

* `ids` - The ServiceAction IDs
* `names` - The ServiceAction names
* `descriptions` - The ServiceAction descriptions
* `types` - The ServiceAction types
