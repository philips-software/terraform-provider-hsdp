---
subcategory: "Identity and Access Management (IAM)"
---

# hsdp_iam_client

Retrieve details of an existing OAuth client

## Example Usage

```hcl
data "hsdp_iam_client" "my_client" {
   name = "MYCLIENT"
   application_id = data.hsdp_iam_appliation.my_app.id
}
```

```hcl
output "my_client_id" {
   value = data.hsdp_iam_client.my_client.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the application to look up
* `application_id` - (Required) the UUID of the application the client belongs to

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `description` - The description of the application
* `global_reference_id` - The global reference ID of the application
* `type` - Either `Public` or `Confidential`
* `client_id` - The client id
* `response_types` - Array. Examples of response types are "code id\_token", "token id\_token", etc.
* `scopes` - Array. List of supported scopes for this client
* `default_scopes` - Array. Default scopes. You do not have to specify these explicitly when requesting a token
* `redirection_uris` - Array of valid RedirectionURIs for this client
* `consent_implied` - Flag when enabled, the resource owner will not be asked for consent during authorization flows.
* `access_token_lifetime` - Lifetime of the access token in seconds
* `refresh_token_lifetime` - Lifetime of the refresh token in seconds
* `id_token_lifetime` - (Optional) Lifetime of the jwt token generated in case openid scope is enabled for the client.
