---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_oauth_client

Provides a resource for managing Connect IoT OAuth2 clients

## Example Usage

The following example creates a Connect IoT OAuth client

```hcl
resource "hsdp_connect_mdm_oauth_client" "testclient" {
  name                = "TESTCLIENT"
  description         = "Test client"
  application_id      = data.hsdp_connect_mdm_application.test_app.id
  global_reference_id = "some-ref-here"
  
  scopes         = [ "cn", "introspect", "email", "profile" ]
  default_scopes = [ "cn", "introspect" ]
  
  redirection_uris = [
    "https://foo.bar/auth",
    "https://testapp.cloud.pcftest.com/auth",
  ]

  response_types = ["code", "code id_token"]
  
  user_client = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the client
* `description` - (Required) The description of the client
* `application_id` - (Required) the application ID (GUID) to attach this client to
* `global_reference_id` - (Required) Reference identifier defined by the provisioning user. This reference Identifier will be carried over to identify the provisioned resource across deployment instances (ClientTest, Production). Invalid Characters:- "[&+’";=?()\[\]<>]
* `scopes` - (Optional) Array. List of supported scopes for this client
* `default_scopes` - (Optional) Array. Default scopes. You do not have to specify these explicitly when requesting a token
* `redirection_uris` - (Optional) Array of valid RedirectionURIs for this client
* `user_client` - (Optional, bool)
* `client_revoked` - (Optional, bool)

~> The `application_id` only accept MDM Application IDs. Using an IAM Proposition ID will not work, even though they might look similar.


## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the client
* `disabled` - True if the client is disabled e.g. because the Org is disabled
* `client_id` -  The client id
* `password` - The password

## Import

An existing client can be imported using `terraform import hsdp_connect_mdm_oauth_client`, e.g.

```shell
terraform import hsdp_iam_client.myclient a-guid
```
