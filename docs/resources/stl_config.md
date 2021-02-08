# hsdp_stl_config
Manage configuration of a STL device

## Example usage
```hcl
data "hsdp_stl_device" "sme100" {
  serial_number = "S4439394855830303"
}

resource "hsdp_stl_config" "sme100" {
  serial_number = data.hsdp_stl_device.sme100.serial_number
  
  firewall_exceptions {
    tcp = ["8080", "4443"]
    udp = ["53"]
  }

  logging {
    raw_config = file(var.raw_fluentbit_config)
    hsdp_product_key = var.logging_product_key
    hsdp_shared_key = var.logging_shared_key
    hsdp_secret_key = var.logging_secret_key
    hsdp_logging_endpoint = var.logging_endpoint
  }

  cert {
    name = "test1"
    private_key = hsdp_pki_cert.test1.private_key_pem
    cert = hsdp_pki_cert.test1.cert_pem
  }
  
  cert {
    name = "test2"
    private_key_pem = hsdp_pki_cert.test2.private_key_pem
    cert_pem = hsdp_pki_cert.test2.cert_pem
  }
}
```


## Argument reference
* `serial_number` - (Required) The serial of the device this config should be applied to
* `firewall_exceptions` - (Optional) Firewall exceptions
  * `tcp` - (Optional, list(string)) Array of TCP ports to allow
  * `udp` - (Optional, list(string)) Array of UDP ports to allow
* `cert` - (Optional) A custom certificate to install on the device
  * `name` - (Required) Name of the certificate
  * `cert_pem`  - (Required) The certificate in PEM format
  * `private_key_pem` - (Required) the private key of the certificate in PEM format  
* `logging` - (Optional) Log forwarding and fluent-bit logging configuration for the device
  * `raw_config` - (Optional) Fluent-bit raw configuration to use
  * `hsdp_product_key` - (Optional) the HSDP logging product key
  * `hsdp_shared_key` - (Optional) the HSDP logging shared key
  * `hsdp_secret_key` - (Optional) the HSDP logging secret key
  * `hsdp_logging_endpoint` - (Optional) The HSDP logging endpoint

## Attribute reference