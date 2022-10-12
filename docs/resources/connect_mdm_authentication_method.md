---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_authentication_method

Create and manage MDM Authentication Method resources

## Example Usage

```hcl
resource "hsdp_connect_mdm_authentication_method" "some_auth_method" {
  name        = "some-authentication-method"
  description = "An authentication method"
  
  login_id = var.login_id
  password = random_password.generated.result
  
  client_id     = var.client_id
  client_secret = var.client_secret
  
  auth_method = "Basic"
  auth_url    = "https://api.login.app.hsdp.io"
  api_version = "3"
}

resource "random_password" "generated" {
  length           = 16
  special          = true
  min_upper        = 1
  min_lower        = 1
  min_numeric      = 1
  min_special      = 1
  override_special = "-!@#.:_?{$"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the device group
* `description` - (Optional) A short description of the device group
* `login_name` - (Required) The login name to use
* `password` - (Required) The password to use
* `client_id` - (Required) The client ID to use
* `client_secret` - (Required) the client secret to use
* `auth_url` - (Required) The authentication URL to use
* `auth_method` - (Required) the authentication method to use [`Bearer` | `Basic`]
* `api_version` - (Required) the API version to use
* `organization_id` - (Optional) The organization ID to associate this method to

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `Group/${GUID}`)
* `guid` - The GUID of the service action
