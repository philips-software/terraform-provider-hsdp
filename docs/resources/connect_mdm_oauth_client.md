---
subcategory: "Master Data Management (MDM)"
page_title: "HSDP: hsdp_connect_mdm_oauth_application"
description: |-
  Manages HSDP Connect MDM OAuth applications
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

  scopes = [
    "?.?.dsc.service.readAny",
    "?.?.prf.profile-custom.UpdateAny",
    "?.*.prf.profile-firmware.UpdateAny",
    "?.?.prf.profile-firmware.UpdateAny",
    "?.?.prf.profile-firmware.UpdateOwn"
  ]
  default_scopes = [
    "?.?.dsc.service.readAny",
    "?.?.prf.profile-custom.UpdateAny",
    "?.*.prf.profile-firmware.UpdateAny",
    "?.?.prf.profile-firmware.UpdateAny",
    "?.?.prf.profile-firmware.UpdateOwn"
  ]

  iam_scopes = [
    "tdr.contract"
  ]
  iam_default_scopes = [
    "tdr.contract"
  ]
  
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
* `iam_scopes` - (Optional) Array. List of supported scopes for this client's IAM instance
* `iam_default_scopes` - (Optional) Array. Default scopes to set for this client's IAM instance
* `bootstrap_client_scopes` - (Optional) Array. List of supported scopes for the bootstrap client
* `bootstrap_client_default_scopes` - (Optional) Array. Default scopes for the bootstrap client. You do not have to specify these explicitly when requesting a token
* `bootstrap_client_iam_scopes` - (Optional) Array. List of supported scopes for this bootstrap client's IAM instance
* `bootstrap_client_iam_default_scopes` - (Optional) Array. Default scopes to set for this bootstrap client's IAM instance
* `redirection_uris` - (Optional) Array of valid RedirectionURIs for this client
* `user_client` - (Optional, bool)
* `client_revoked` - (Optional, bool)

~> The `application_id` only accept MDM Application IDs. Using an IAM Proposition ID will not work, even though they might look similar.
~> If `user_client` is false, only `scopes`, `default_scopes`, `iam_scopes` and `iam_default_scopes` are allowed.

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the client
* `disabled` - True if the client is disabled e.g. because the Org is disabled
* `client_id` -  The client id
* `password` - The password
* `bootstrap_client_id` - The bootstrap client ID
* `bootstrap_client_secret` - The boostrap client secret
* `bootstrap_client_guid_system` - The external system bootstrap client associated with resource (this would point to an IAM deployment)
* `bootstrap_client_guid_value` - The external value of the bootstrap client associated with this resource (this would be an underlying IAM OAuth2 client GUID)
* `client_guid_system` - The external system client associated with resource (this would point to an IAM deployment)
* `client_guid_value` - The external value client associated with this resource (this would be an underlying IAM OAuth2 client GUID)

## Import

An existing client can be imported using `terraform import hsdp_connect_mdm_oauth_client`, e.g.

```shell
terraform import hsdp_iam_client.myclient a-guid
```
