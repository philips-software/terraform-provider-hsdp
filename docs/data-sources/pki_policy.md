# hsdp_pki_policy
Retrieves the HSDP PKI Policy CA and CRL

# Example Usage

```hcl

data "hsdp_pki_policy" "info" {
}

output "policy_ca" {
  value = data.hsdp_pki_policy.info.ca_pem
}
```
# Argument reference
* `region` - (Optional) the HSDP PKI regional selection
* `environment` - (Optional) the HSDP PKI environment to use [`client_test` | `prod`]

# Attribute reference

* `ca_pem` - The root CA in PEM format
* `crl_pem` - The root CRL in PEM format
