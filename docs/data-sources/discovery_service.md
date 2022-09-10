---
subcategory: "Discovery"
---

# hsdp_discovery_service

Discover standard services details

## Example usage

```hcl
// Discover by name
data "hsdp_discovery_service" "my_service" {
    name = "My Service"
}

// Discover by tag
data "hsdp_discovery_service" "another_service" {
    tag = "another-service"
}
```

## Argument Reference

One of the below arguments should be provided

* `name` - (Optional) The name of the Service to discover
* `tag` - (Optional) The tag of the Service to discover

~> The current version of the Discovery service only supports `User` and `Device` principals. `Service identities` are not supported (yet)

### Principal

You can optionally specify a principal (`User` or `Device`) to use for performing the service discovery calls

* `principal` - (Optional) The optional principal to use for this resource
  * `username` - (Optional) The username of the user or device
  * `password` - (Optional) The password of the user or device
  * `oauth2_client_id` - (Optional) The MDM OAuth2 client ID to use for token exchange
  * `oauth2_password` - (Optional) The MDM OAuth2 client password to use for token exchange

~> An `MDM OAuth2 client` with the `?.?.dsc.service.readAny` set should be used for retrieving the principal token. At the time of writing this document (September 2022) you will almost certainly require a `principal` block for correct operation of this data source

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the service
* `resource_type` - The resource type of the service
* `is_trusted` - Wether this service is a trusted one
* `actions` - (list(string)) A list of actions supported by the service
* `urls` - (list(string)) The list of URLs of this service. Ordered is significant
