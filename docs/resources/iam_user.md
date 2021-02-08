# hsdp_iam_user
Provides a resource for managing HSDP IAM application under a proposition. When a new user is created an invitation email is triggered with a validity of 24 hours. If the user fails to activate his/her account within this time period the account will be removed in the backend and the resource should be removed or cleared from tfstate.

>Typically this resource is used to only test account. We highly recommend using the IAM Self serviceUI which HSDP provides for day to day user management tasks


## Example Usage

The following example creates a user. 

```hcl
resource "hsdp_iam_user" "developer" {
  login           = "developer"
  email           = "developer@1e100.io"
  first_name      = "Devel"
  last_name       = "Oper"
  mobile          = "316123456789"
  organization_id = hsdp_iam_org.testdev.id
}
```

## Argument Reference

The following arguments are supported:

* `login` - (Required) The login ID of the user (NEW since v0.4.0)
* `email` - (Required) The email address of the user
* `first_name` - (Required) First name of the user
* `last_name` - (Required) Last name of the user
* `mobile` - (Required) Mobile number of the user. E.164 format
* `organization_id` - (Required) The managing organization of the user

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the user

## Import

An existing user can be imported using `terraform import hsdp_iam_user`, e.g.

```shell
> terraform import hsdp_iam_user.developer a-guid
```

