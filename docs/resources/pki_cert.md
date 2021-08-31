# hsdp_pki_cert

Create and manage HSDP PKI leaf certificates

## Example usage

```hcl
resource "hsdp_pki_cert" "cert" {
  tenant_id = hsdp_pki_tenant.tenant.id
  role      = "ec384"
  
  common_name = "myapp.com"
  alt_names    = "myapp.io,www.myapp.io"
  ip_sans     = []
  uri_sans    = []
  other_sans  = []
  ttl         = "720h"
  
  exclude_cn_from_sans = false
}
```

## Argument reference

* `tenant_id` - (Required) The tenant ID to create this certificate under
* `role` - (Required) the Role to use as defined under a PKI Tenant resource
* `common_name` - (Required) The common name to use
* `alt_names` - (Optional) Alternative names to use, comma separated list.
* `ip_sans` - (Optional, list(string)) A list of IP SANS to include
* `uri_sans` - (Optional, list(string)) A list of URI SANS to include
* `other_sans` - (Optional, list(string)) A list of other SANS to include
* `ttl` - (Optional, string regex `[0-9]+[hms]$`) The TTL, example `720h` for 1 month
* `exclude_cn_from_sans` - (Optional) Exclude common name from SAN

## Attribute reference

* `cert_pem` - The certificate in PEM format
* `private_key_pem` - The private key in PEM format
* `issuing_ca_pem` - The issuing CA certicate in PEM format
* `serial_number` - The certificate serial number (equal to resource ID)
* `expiration` - (int) The Unix timestamp when the certificate will expire
* `ca_chain_pem` - The full CA chain in PEM format

## Importing

Importing a HSDP PKI certificate is supported but not recommended as the private key will be missing,
rendering the resource more or less useless in most cases. You can import a certificate using the serial number
