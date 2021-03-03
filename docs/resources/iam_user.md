# hsdp_iam_user
Provides a resource for managing an HSDP IAM user. 
When a new user is created an invitation email is triggered with a validity of 72 hours. 
If not activated within this period IAM will purge the account.
The provider will recreate the user in that case.

## Example Usage

The following example creates a user. 

```hcl
resource "hsdp_iam_user" "developer" {
  login           = "developer"
  email           = "developer@1e100.io"
  first_name      = "Devel"
  last_name       = "Oper"
  organization_id = hsdp_iam_org.testdev.id
}
```

## Argument Reference

The following arguments are supported:

* `organization_id` - (Required) The managing organization of the user

* `login` - (Required) The login ID of the user (NEW since v0.4.0)
* `email` - (Required) The email address of the user
* `first_name` - (Required) First name of the user
* `last_name` - (Required) Last name of the user
* `mobile` - (Optional) Mobile number of the user. E.164 format
* `password` - (Optional) When specified this will skip the email activation 
  flow and immediately activate the IAM account. **Very Important**: you are responsible
  for sharing this password with the new IAM user through some channel of communication. 
  No email will be triggered by the system. If unsure, do not set a password so the normal 
  email activation flow is followed. Finally, any password value changes after user creation
  will have no effect on the users' actual password.

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the user

## Import

An existing user can be imported using `terraform import hsdp_iam_user`, e.g.

```shell
> terraform import hsdp_iam_user.developer a-guid
```

