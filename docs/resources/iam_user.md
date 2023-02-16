---
subcategory: "Identity and Access Management (IAM)"
page_title: "HSDP: hsdp_iam_user"
description: |-
  Manages HSDP IAM Users
---

# hsdp_iam_user

Provides a resource for managing a HSDP IAM user. When a new user is created an invitation email is triggered with a validity of 72 hours.
In case the user hasn't activated their account within this 72-hour period you can use the [hsdp_iam_activation_email](https://registry.terraform.io/providers/philips-software/hsdp/latest/docs/resources/iam_activation_email)
resource to resend the email. Identifying unactivated users can be done using the [hsdp_iam_users](https://registry.terraform.io/providers/philips-software/hsdp/latest/docs/data-sources/iam_users) data source.

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
* `email` - (Semi-Required) The email address of the user
* `first_name` - (Required) First name of the user
* `last_name` - (Required) Last name of the user
* `mobile` - (Optional) Mobile number of the user. E.164 format
* `password` - (Optional) When specified this will skip the email activation
  flow and immediately activate the IAM account. **Very Important**: you are responsible
  for sharing this password with the new IAM user through some channel of communication.
  No email will be triggered by the system. If unsure, do not set a password so the normal
  email activation flow is followed. Finally, any password value changes after user creation
  will have no effect on the users' actual password.
* `preferred_language` - (Optional) Language preference for all communications.
  Value can be a two letter language code as defined by ISO 639-1 (en, de) or it can be a combination
  of language code and country code (en-gb, en-us). The country code is as per ISO 3166 two letter code (alpha-2)
* `preferred_communication_channel` - (Optional) Preferred communication channel.
  Email and SMS are supported channels. Email is the default channel if e-mail address is provided.
  Values supported: [ `email` | `sms` ]

> Use the `preferred_*` arguments sparingly as they will reset values if the user has changed these outside of Terraform

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the user

## Import

An existing user can be imported using `terraform import hsdp_iam_user`, e.g.

```shell
> terraform import hsdp_iam_user.developer a-guid
```
