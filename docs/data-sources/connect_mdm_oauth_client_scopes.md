---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_proposition

Retrieve details of an existing proposition

## Example Usage

```hcl
data "hsdp_connect_mdm_oauth_client_scopes" "all" {
}
```

```hcl
output "mdm_oauth_client_scope_names" {
   value = data.hsdp_connect_mdm_oauth_client_scopes.all.ids
}
```

## Attributes Reference

The following attributes are exported:

* `ids` - The client scope IDs
* `names` - the client socpe names
* `actions` - The client scope actions
* `propositions` - The client scope actions
* `organizations` - The client organization list
* `services` - The client scope services list
* `bootstrap_enabled` - Bootstrap enabled list
