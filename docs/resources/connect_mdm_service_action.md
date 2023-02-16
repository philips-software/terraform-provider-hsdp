---
subcategory: "Master Data Management (MDM)"
page_title: "HSDP: hsdp_connect_mdm_service_action"
description: |-
  Manages HSDP Connect MDM Service actions
---

# hsdp_connect_mdm_service_action

Create and manage MDM ServiceAction resources

## Example Usage

```hcl
resource "hsdp_connect_mdm_service_action" "some_action" {
  name        = "Some action"
  description = "A tenant service action that does something"
  
  standard_service_id = hsdp_connect_mdm_standard_service.some_service.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service action
* `description` - (Optional) A short description of the service action
* `standard_service_id` - (Required) Reference to a Standard Service

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `ServiceAction/${GUID}`)
* `guid` - The GUID of the service action
