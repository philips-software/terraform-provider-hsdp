# hsdp_pki_tenant

Onboard tenant to PKI Service. Cloud foundry users with SpaceDeveloper role can onboard tenant

> This resource is only available when `uaa_*` (Cloud foundry) and `iam` credentials are set

# Example usage

```hcl
resource "hsdp_pki_tenant" "tenant" {
  organization_name = "client-my-org"
  space_name = "prod"
  
  iam_orgs = [
    var.iam_org_id
  ]
  
  ca {
    common_name = "Common Name Here"
  }
  
  role {
    name = "role a"
    allow_any_name = true
    allow_ip_sans = true
    allow_subdomains = true
    allowed_domains = []
    allowed_other_sans = []
    allowed_uri_sans = []
    client_flag = true
    server_flag = true
    enforce_hostnames = false
    key_bits = 384
    key_type = "ec"
  }
}
```