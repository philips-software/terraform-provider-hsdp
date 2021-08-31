# hsdp_edge_sync
Synchronizes device configuration. This resource can be used to batch sync requests
of a device e.g. you can add all resource configs to the trigger hash and disable sync
per resource to batch syncs down to a single one as part of the `apply` stage.

## Argument reference
* `serial_number` - (Required) Serial number of the device to sync
* `triggers` - (Required, Hashmap) Create dependencies on other resources
