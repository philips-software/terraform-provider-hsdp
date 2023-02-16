---
subcategory: "HealthSuite Edge"
page_title: "HSDP: hsdp_edge_sync"
description: |-
  Manages HSDP Edge synchronizations
---

# hsdp_edge_sync

Synchronizes device configuration. This resource can be used to batch sync requests
of a device e.g. you can add all resource configs to the trigger hash and disable sync
per resource to batch syncs down to a single one as part of the `apply` stage.

## Argument reference

* `serial_number` - (Required) Serial number of the device to sync
* `triggers` - (Required, Hashmap) Create dependencies on other resources
* `principal` - (Optional) The optional principal to use for this resource
  * `uaa_username` - (Optional) The UAA username to use
  * `uaa_password` - (Optional) The UAA password to use
  * `region` - (Optional) Region to use. When not set, the provider config is used
  * `endpoint` - (Optional) The endpoint URL to use if applicable. When not set, the provider config is used
