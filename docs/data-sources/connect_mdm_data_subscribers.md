---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_data_subscribers

Retrieve details of global DataSubscribers

## Example Usage

```hcl
data "hsdp_connect_mdm_data_subscribers" "all" {
}
```

```hcl
output "data_subscribers_names" {
   value = data.hsdp_connect_data_subscribers.all.names
}
```

## Attributes Reference

The following attributes are exported:

* `ids` - The DataSubscriber IDs
* `names` - The DataSubscriber names
* `configurations` - The configurations
* `subscriber_guids` - The subscriber GUIDs
* `subscriber_type_ids` - The subscriber type IDs
