---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_client_scopes

Retrieve client scope dictionary

## Example Usage

```hcl
data "hsdp_connect_mdm_oauth_client_scopes" "all" {
}
```

```hcl
output "mdm_oauth_client_scopes" {
   value = data.hsdp_connect_mdm_oauth_client_scopes.all.scopes
}
```

## Attributes Reference

The following attributes are exported:

* `ids` - The client scope IDs
* `names` - the client dictionary entry names
* `scopes` - The effective scopes
* `actions` - The client scope actions
* `propositions` - The client scope actions
* `organizations` - The client organization list
* `services` - The client scope services list
* `bootstrap_enabled` - Bootstrap enabled list
