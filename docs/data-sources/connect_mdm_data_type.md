---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_data_type

Retrieve information on MDM DataType resource

## Example Usage

```hcl
data "hsdp_connect_mdm_data_type" "some_type" {
  name           = "tf-some-data-type"
  proposition_id = data.hsdp_connect_mdm_proposition.first.id
}

output "data_type_guid" {
  value = data.hsdp_connect_mdm_data_type.some_type.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the device group
* `proposition_id` - (Required) The proposition to which the data type is associated with

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `DataType/${GUID}`)
* `guid` - The GUID of the service action
* `description` - A short description of the device group
* `tags` - (list(string)) Tags associated with this data type
