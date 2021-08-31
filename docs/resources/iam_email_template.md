# hsdp_iam_email_template

There are certain IAM flows that trigger email notifications to the user. The email delivered to users will use the IAM specific email template by default.
This resource allows you to manage template for your HSDP IAM organization. [Further reading](https://www.hsdp.io/documentation/identity-and-access-management-iam/how-to/customize#_email_template_customization).

## Types of templates

The various template types supported by IAM are:

| Type | Description |
|------|--------------|
| ACCOUNT_ALREADY_VERIFIED | If a user’s account is already verified and activated and the resend verification is triggered, the user gets this email with a message that the account is already verified. |
| ACCOUNT_UNLOCKED | If a user’s account is unlocked by an administrator, the user gets this email notification. |
| ACCOUNT_VERIFICATION | When a user gets registered within an organization, an account verification email will be sent to the user. The email message will contain an account verification link that will redirect users to the set password page through which the user can set a password and complete the registration process. |
| MFA_DISABLED | If multi-factor authentication is disabled for a user, the user will get this email notification. |
| MFA_ENABLED | If multi-factor authentication is enabled for a user, the user will get this email notification. |
| PASSWORD_CHANGED | If a user’s password is changed successfully, the user will get this email notification. |
| PASSWORD_EXPIRY | If a user’s password is about to expire, this email will be sent to the user with a link to change the password. |
| PASSWORD_FAILED_ATTEMPTS | If there are multiple attempts to change a user account password with invalid current password, then the user will get this email notification warning the user about malicious login attempts. This notification will be sent after 5 invalid attempts. |
| PASSWORD_RECOVERY | If a user triggers the forgot password flow, a password reset email message will be sent to the user. The email message will contain a reset password link that will redirect the user to the reset password page, through which the user can set a new password. |

## Placeholders

Email template supports adding certain placeholders and redirect link in message body based on template type specified. This allows client to templatize certain parts of email template with dynamic data based on whom the email is targeted.

All template types support the following placeholders in subject and message

* `{{user.email}}` - Email address of user
* `{{user.userName}}` - Unique login ID of the user
* `{{user.givenName}}` - User's first name
* `{{user.familyName}}` - User's last name

Various template types supported by the platform are

### ACCOUNT_VERIFICATION

When a user gets registered within an organization, an account verification email will be sent to user. Clicking the account verification link, will activate the user account.
The following placeholders are supported in this template

* `{{link.verification}}` - Account verification uri
* `{{template.linkExpiryPeriod}}` - How long the verification link is valid (in hours)
If link field is configured to a custom URL, upon clicking the verification link, the user will be redirected to the specified link with an OTP appended to it. If no link is set, HSDP default link will be used for account activation.

### ACCOUNT_ALREADY_VERIFIED

If a user's account is already verified and activated, and if resend verification is triggered, user gets this mail with a message that account is already verified. No link needs to be configured for this template.

### ACCOUNT_UNLOCKED

If a user's account is unlocked by an administrator, user gets this email notification. No link needs to be configured for this template.

### PASSWORD_RECOVERY

If a user triggers forgot password flow, a password reset email message will be sent to user. The email message will contain a reset password link that will redirect the user to reset password page through which user can set a new password.
The following placeholders are supported in this template

* `{{link.passwordReset}}` - Reset password uri
The link can be configured to a custom page where the user can set new password. If no link is set, HSDP default link will be used for password reset.

### PASSWORD_EXPIRY

If a user password is about to expire, this email will be sent to user with a link to change the password.
The following placeholders are supported in this template

* `{{link.passwordChange}}` - Change password uri
* `{{password.expiresAfterPeriod}}` - How long the current password is valid (in days)
The link can be configured to a custom login page. Upon login the user can be redirected to change password page where the user should be asked to enter old password and new password. If no link is set, HSDP default link will be used for password change.

### PASSWORD_FAILED_ATTEMPTS

If there are multiple attempts to change user account password with invalid current password, then the user will get this email notification warning user about malicious login attempts. This notification will be sent after 5 invalid attempts.
The following placeholders are supported in this template

* `{{user.lockoutPeriod}}` - How long the account will be in locked state (in minutes)

### PASSWORD_CHANGED

If user's password is changed successfully, user gets this email notification. No link need to be configured for this template.

### MFA_ENABLED

If multi-factor authentication is enabled for a user, user will get this email notification. No link need to be configured for this template.

### MFA_DISABLED

If multi-factor authentication is disabled for a user, user will get this email notification. No link need to be configured for this template.

## Example Usage

The following example manages an email template for an org

```hcl
resource "hsdp_iam_email_template" "password_changed" {
  type = "PASSWORD_CHANGED"

  managing_organization = data.hsdp_iam_org.myorg.id
  format                = "HTML"

  subject = "Your IAM account password was changed"
  message = <<EOF
Dear {{user.givenName}},

Your password was recently changed. If this was not initiated
by you please contact support immediately.

Kind regards,
IAM Team
EOF
}

resource "hsdp_iam_email_template" "password_expiry" {
  type = "PASSWORD_EXPIRY"

  managing_organization = data.hsdp_iam_org.myorg.id
  format                = "HTML"

  subject = "Your IAM account password was changed"
  message = <<EOF
Dear {{user.givenName}},

Your password will expire in {{password.expiresAfterPeriod}} day(s). Please set a new
password using the below link:

{{link.passwordChange}}

Kind regards,
IAM Team
EOF

}
```

## Argument Reference

The following arguments are supported:

* `managing_organization` - (Required) The UUID of the IAM Org to apply this email template to
* `type` - (Required) The email template. See the `Type` table above for available values
* `format` - (Required) The template format. Must be `HTML` currently
* `message` - (Required) A boolean value indicating if challenges are enabled at organization level. If the value is set to true, `challenge_policy` attribute is mandatory.
* `locale` - (Optional) The locale of the template. When not specified the template will become the default. Only a single default template is allowed of course.
* `from` - (Optional) The From field of the email. Default value is `default`
* `subject` - (Optional) The Subject line of the email. Default value is `default`
* `link` - (Optional) A clickable link, depends on the template `type`

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the email template

## Import

Importing is supported but not recommended as the `message` body is not returned when reading out a template via the IAM API
