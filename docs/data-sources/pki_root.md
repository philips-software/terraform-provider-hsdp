# hsdp_pki_root
Retrieves the HSDP PKI Root CA and CRL

# Example Usage

```hcl

data "hsdp_pki_root" "info" {
}

output "root_ca" {
  value = data.hsdp_pki_root.info.ca_pem
}

output "root_crl" {
  value = data.hsdp_pki_root.info.crl_pem
}
```
# Argument reference
* `region` - (Optional) the HSDP PKI regional selection
* `environment` - (Optional) the HSDP PKI environment to use [`client_test` | `prod`]
 
# Attribute reference

* `ca_pem` - The root CA in PEM format
* `crl_pem` - The root CRL in PEM format
