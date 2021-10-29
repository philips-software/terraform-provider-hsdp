---
subcategory: "HealthSuite Edge"
---

# hsdp_edge_device

Represents an Edge device

## Example usage

```hcl
data "hsdp_edge_device" "sme100" {
  serial_number = "S4439394855830303"
}

output "sme100_status" {
  value = data.hsdp_edge_device.sme100.status
}
```

## Argument reference

* `serial_number` - (Required) the serial number of the device

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - The device ID
* `region` - The region to which this device is connected to
* `state` - State of the device
* `primary_interface_ip` - The IP of the primary interface
