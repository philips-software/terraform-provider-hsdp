---
subcategory: "Identity and Access Management (IAM)"
---

# hsdp_iam_introspect

Introspects the ORG admin account in use by the provider

~> This data source only works if the provider is configured with OAuth2 client credentials (`oauth2_client_id` and `oauth2_password`)

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

### Principal

You can optionally specify a principal (service identity, user or device) to use for performing the introspect call

* `principal` - (Optional) The optional principal to use for this resource
  * `service_id` - (Optional) The IAM service ID
  * `service_private_key` - (Optional) The IAM service private key to use
  * `username` - (Optional) The username of the user or device
  * `password` - (Optional) The password of the user or device
  * `oauth2_client_id` - (Optional) The OAuth2 client ID to use for token exchange
  * `oauth2_password` - (Optional) The OAuth2 client password to use for token exchange
  * `region` - (Optional) Region to use. When not set, the provider config is used
  * `environment` - (Optional) Environment to use. When not set, the provider config is used
  * `endpoint` - (Optional) The endpoint URL to use if applicable. When not set, the provider config is used

## Attributes Reference

The following attributes are exported:

* `managing_organization` - The managing organization of the Org admin user
* `subject` - The subject of the token, as defined in JWT RFC7519.
  Usually a machine-readable identifier of the resource owner who authorized this token.
* `issuer` - String representing the issuer of this token, as defined in JWT
* `username` - The username (email) of the Org admin user
* `token` - The current session token
* `effective_permissions` - When an `organization_context` GUID is provided this
  contains the list of effective permissions
* `token_type` - The type of token
* `identity_type` - The identity type, example: `Service`
* `scopes` - The list of scopes associated with the token
