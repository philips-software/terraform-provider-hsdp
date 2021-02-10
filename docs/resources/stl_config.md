# hsdp_stl_config
Manage configuration of a STL device. Set `sync` to true to immediately sync the config to the device, otherwise
you should create a dependency on a `hsdp_stl_sync` resource to batch sync changes.

## Example usage
```hcl
data "hsdp_stl_device" "sme100" {
  serial_number = "S4439394855830303"
}

resource "hsdp_stl_config" "sme100" {
  serial_number = data.hsdp_stl_device.sme100.serial_number
  
  firewall_exceptions {
    tcp = [8080, 4443]
    udp = [53]
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
  * `tcp` - (Optional, list(int)) Array of TCP ports to allow
  * `udp` - (Optional, list(int)) Array of UDP ports to allow
* `logging` - (Optional) Log forwarding and fluent-bit logging configuration for the device
  * `raw_config` - (Optional) Fluent-bit raw configuration to use
  * `hsdp_logging` - (Optiona, boolean) Enable or disable HSDP Logging feature   
  * `hsdp_product_key` - (Optional) the HSDP logging product key
  * `hsdp_shared_key` - (Optional) the HSDP logging shared key
  * `hsdp_secret_key` - (Optional) the HSDP logging secret key
  * `hsdp_ingestor_host` - (Optional) The HSDP logging endpoint

## Attribute reference
* `last_update` - RFC3339 timestamp of last update. Can be used to trigger `hsdp_stl_sync`
