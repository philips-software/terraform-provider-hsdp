---
subcategory: "Identity and Access Management"
---

# hsdp_iam_introspect

Introspects the ORG admin account in use by the provider

## Example Usage

```hcl
data "hsdp_iam_introspect" "admin" {}
```

```hcl
output "admins_org" {
   value = data.hsdp_iam_introspect.admin.managing_organization
}
```

## Attributes Reference

The following attributes are exported:

* `managing_organization` - The managing organization of the Org admin user
* `username` - The username (email) of the Org admin user
* `token` - The current session token
