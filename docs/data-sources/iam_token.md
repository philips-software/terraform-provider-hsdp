---
subcategory: "Identity and Access Management (IAM)"
---

# hsdp_iam_token

Retrieve IAM tokens for use in other resources. The provider
configured IAM entity credentials are used to generate these tokens.

~> This data source regenerates the tokens each time a new plan is created.

## Example Usage

```hcl
data "hsdp_iam_token" "iam" {
}
```

```hcl
output "access_token" {
   value = data.hsdp_iam_token.iam.access_token
}
```

## Attributes Reference

The following attributes are exported:

* `access_token` - (string) An IAM Access token. This has a limited TTL, usually 30 minutes.
* `id_token` - (string) An IAM ID token. This has a limited TTL, usually 30 minutes.
* `expires_at` - (number) The Unix timestamp when the access token expires
