# hsdp_pki_cert

# Example usage

```hcl
resource "hsdp_pki_cert" "cert" {
  tenant_id = hsdp_pki_tenant.tenant.id
  role      = "ec384"
  
  common_name = "myapp.com"
  alt_name    = "myapp.io"
  ip_sans     = ""
  uri_sans    = ""
  other_sans  = ""
  ttl         = 86400
  format      = "ec"
  
  exclude_cn_from_sans = false
}
```