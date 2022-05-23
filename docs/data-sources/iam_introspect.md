---
subcategory: "Identity and Access Management (IAM)"
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

## Argument Reference

* `token` - (Optional) the token to introspect. Uses default token otherwise
* `organization_context` - (Optional) Does a contextual introspect the IAM Organization associated
   with the GUID. The `effective_permissions` attribute will contain the list of permissions.

## Attributes Reference

The following attributes are exported:

* `managing_organization` - The managing organization of the Org admin user
* `username` - The username (email) of the Org admin user
* `token` - The current session token
* `effective_permissions` - When an `organization_context` GUID is provided this
  contains the list of effective permissions
* `token_type` - The type of token
* `identity_type` - The identity type, example: `Service`
