# hsdp_stl_custom_cert
Manage custom certificates for STL devices

## Example usage
```hcl
resource "hsdp_stl_custom_cert" "cert" {
  serial_number = var.serial_number
  
  name = "terrakube.com"
  cert_pem = var.cert_pem
  private_key_pem = var.private_key_pme
}
```
## Argument reference
* `serial_number` - (Required) Device to attach the cert to
* `name` - (Required) Name of the certificate
* `cert_pem`  - (Required) The certificate in PEM format
* `private_key_pem` - (Required) the private key of the certificate in PEM format
* `sync` - (Optional) When set to true syncs config automatically after a mutation. Default is false

## Attribute reference
* `id` - The id of the custom certificate

## Importing
Importing a custom certificate is supported but not recommended.