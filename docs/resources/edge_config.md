---
subcategory: "HealthSuite Edge"
page_title: "HSDP: hsdp_edge_config"
description: |-
  Manages HSDP Edge configurations
---

# hsdp_edge_config

Manage configuration of an Edge device. Set `sync` to true to immediately sync the config to the device, otherwise
you should create a dependency on a `hsdp_edge_sync` resource to batch sync changes.

## Example usage

```hcl
data "hsdp_edge_device" "sme100" {
  serial_number = "S4439394855830303"
}

resource "hsdp_edge_config" "sme100" {
  serial_number = data.hsdp_edge_device.sme100.serial_number
  
  firewall_exceptions {
    ensure_tcp = [2575]
    ensure_udp = [2345]
  }

  logging {
    raw_config = file(var.raw_fluentbit_config)

    hsdp_logging = true
    hsdp_product_key = var.logging_product_key
    hsdp_shared_key = var.logging_shared_key
    hsdp_secret_key = var.logging_secret_key
    hsdp_ingestor_host = var.logging_endpoint
  }
}
```

## Argument reference

* `serial_number` - (Required) The serial of the device this config should be applied to
* `firewall_exceptions` - (Optional) Firewall exceptions
  * `ensure_tcp` - (Optional, list(int)) Array of TCP ports to add. Conflicts with `tcp`
  * `ensure_udp` - (Optional, list(int)) Array of UDP ports to add. Conflicts with `udp`
  * `tcp` - (Optional, list(int)) Array of TCP ports to allow. Conflicts with `ensure_tcp`
  * `udp` - (Optional, list(int)) Array of UDP ports to allow. Conflicts with `ensure_udp`
  * `clear_on_destroy` - (Optional, boolean) When set to true, clears specified ports on destroy. Default is `true`
* `logging` - (Optional) Log forwarding and fluent-bit logging configuration for the device
  * `raw_config` - (Optional) Fluent-bit raw configuration to use
  * `hsdp_logging` - (Optional, boolean) Enable or disable HSDP Logging feature
  * `hsdp_product_key` - (Optional) the HSDP logging product key
  * `hsdp_shared_key` - (Optional) the HSDP logging shared key
  * `hsdp_secret_key` - (Optional) the HSDP logging secret key
  * `hsdp_ingestor_host` - (Optional) The HSDP logging endpoint
* `sync` (Optional, boolean) - When set to true syncs the config after mutations. Default is true.
  Set this to false if you want to batch sync to your device using `hsdp_edge_sync`
* `principal` - (Optional) The optional principal to use for this resource
  * `uaa_username` - (Optional) The UAA username to use
  * `uaa_password` - (Optional) The UAA password to use
  * `region` - (Optional) Region to use. When not set, the provider config is used
  * `endpoint` - (Optional) The endpoint URL to use if applicable. When not set, the provider config is used
