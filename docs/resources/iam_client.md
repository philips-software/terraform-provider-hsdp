# hsdp_iam_client

Provides a resource for managing HSDP IAM client of an application under a proposition.

## Example Usage

The following example creates a client

```hcl
resource "hsdp_iam_client" "testclient" {
  name                = "TESTCLIENT"
  description         = "Test client"
  type                = "Public"
  client_id           = "testclient"
  password            = "Password@123"
  application_id      = hsdp_iam_application.testtapp.id
  global_reference_id = "some-ref-here"
  
  scopes = [ "cn", "introspect", "email", "profile" ]

  default_scopes = [ "cn", "introspect" ]


  redirection_uris = [
    "https://foo.bar/auth",
    "https://testapp.cloud.pcftest.com/auth",
  ]

  response_types = ["code", "code id_token"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the client
* `description` - (Required) The description of the client
* `type` - (Required) Either `Public` or `Confidential`
* `client_id` - (Required) The client id 
* `password` - (Required) The password to use (8-16 chars, at least one capital, number, special char)
* `application_id` - (Required) the application ID (GUID) to attach this client to
* `global_reference_id` - (Required) Reference identifier defined by the provisioning user. This reference Identifier will be carried over to identify the provisioned resource across deployment instances (ClientTest, Production). Invalid Characters:- "[&+’";=?()\[\]<>]
* `response_types` - (Required) Array. Examples of response types are "code id\_token", "token id\_token", etc.
* `scopes` - (Required) Array. List of supported scopes for this client
* `default_scopes` - (Required) Array. Default scopes. You do not have to specify these explicitly when requesting a token
* `redirection_uris` - (Required) Array of valid RedirectionURIs for this client
* `consent_implied` - (Optional) Flag when enabled, the resource owner will not be asked for consent during authorization flows.
* `access_token_lifetime` - (Optional) Lifetime of the access token in seconds. If not specified, system default life time (1800 secs) will be considered.
* `refresh_token_lifetime` - (Optional) Lifetime of the refresh token in seconds. If not specified, system default life time (2592000 secs) will be considered.
* `id_token_lifetime` - (Optional) Lifetime of the jwt token generated in case openid scope is enabled for the client. If not specified, system default life time (3600 secs) will be considered.

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the client
* `disabled` - True if the client is disabled e.g. because the Org is disabled


## Import

An existing client can be imported using `terraform import hsdp_iam_client`, e.g.

```shell
> terraform import hsdp_iam_client.myclient a-guid
```
