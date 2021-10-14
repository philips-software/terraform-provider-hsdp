# hsdp_iam_mfa_policy

Provides a resource for managing HSDP IAM MFA (Multi Factor Authentication) policies

## Example Usage

The following example creates a MFA Policy for an organization

```hcl
resource "hsdp_iam_mfa_policy" "mymfapolicy" {
  type = "SOFT_OTP"
  organization = var.my_org.id
  
  active = true
}
```

And the example below creates a server OTP MFA Policy for an individual user

```hcl
resource "hsdp_iam_mfa_policy" "joes_policy" {
  type = "SERVER_OTP"
  user = var.user_joe.id
  
  active = true
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) the OTP type. Valid values: [`SOFT_OTP` | `SERVER_OTP` | `SERVER_OTP_EMAIL` | `SERVER_OTP_SMS` | `SERVER_OTP_ANY` ]
* `user` - (Required) The user UUID to attach this policy to. Conflicts with `organization`
* `organization` - (Required) The organization to attach this policy to. Conflicts with `user`
* `active` - (Required) Defaults to true. Setting to false will disable MFA for the subject.
* `name` - (Optional) The name of the policy
* `description` - (Optional) The description of the policy

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the MFA policy
* `version` - The version of the MFA policy

## Import

An existing MFA policy can be imported using `terraform import hsdp_iam_mfa_policy`, e.g.

```shell
terraform import hsdp_iam_mfa_policy.mymfapolicy a-guid
```
