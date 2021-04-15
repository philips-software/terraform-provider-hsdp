# hsdp_iam_service
Provides a resource for managing HSDP IAM services of an application under a proposition.

## Example Usage

The following example creates a service

```hcl
resource "hsdp_iam_service" "testservice" {
  name                = "TESTSERVICE"
  description         = "Test service"
  application_id      = var.app_id

  validity            = 12

  scopes              = ["openid"]
  default_scopes      = ["openid"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service
* `description` - (Required) The description of the service
* `application_id` - (Required) the application ID (GUID) to attach this service to
* `scopes` - (Required) Array. List of supported scopes for this service. Minimum: ["openid"]
* `validity` - (Optional) Integer. Validity of service (in months). Minimum: 1, Maximum: 600, Default: 12
* `default_scopes` - (Required) Array. Default scopes. You do not have to specify these explicitly when requesting a token. Minimum: ["openid"]
* `self_private_key` - (Optional) RSA private key in PEM format. When provided, overrides the generated certificate / private key combination of the
  IAM service. This gives you full control over the credentials.
* `self_expires_on` - (Optional) Sets the certificate validity generated from `self_private_key`. When not specified, the generated certificate will have a validity of 5 years.  

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the client
* `service_id` - (Generated) The service id 
* `private_key` - (Generated) The private key of the service
* `expires_on` - (Generated) Date when this service expires
* `organization_id` - The organization ID this service belongs to (via application and proposition)

## Import

Existing services can be imported but they will be missing their private key rendering them pretty much useless. Therefore, we recommend creating them using the provider.
