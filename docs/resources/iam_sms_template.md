# hsdp_iam_sms_template

This resource allows you to provision and manage custom SMS template 
types for an organization. The SMS templates can be registered for 
different locales as well.

## Types of templates

The various template types supported by IAM are:

| Type | Description |
|------|--------------|
| PHONE_VERIFICATION | Send when the users' phone needs to be verified |
| PASSWORD_RECOVERY  | If a user triggers forgot password flow, OTP to reset password will be sent to user. The following placeholders are supported in this template |
| PASSWORD_FAILED_ATTEMPTS | If there are multiple attempts to change user account password with invalid current password, then the user will get this SMS notification warning user about malicious login attempts. This notification will be sent after 5 invalid attempts. The following placeholders are supported in this template |
| MFA_OTP | This SMS template is used for login using OTP for multi-fator authentication. |

## Placeholders

SMS template supports adding certain placeholders in message body based on template type specified. This allows client to templatize certain parts of SMS template with dynamic data based on whom the SMS is targeted.

All template types support the following place holders in message

* `{{user.userName}}` - Unique login ID of the user
* `{{user.givenName}}` - User's first name
* `{{user.familyName}}` - User's last name
* `{{user.displayName}}` - User's display name

### PHONE_VERIFICATION

The following placeholders are supported in this template

* `{{template.otp}}` - Generated OTP.
* `{{template.otpExpiryPeriod}}` - How long the OTP is valid (in minutes)
* `{{template.phoneNumber}}` - Phone number for verification.

### PASSWORD_RECOVERY 

If a user triggers forgot password flow, OTP to reset password will be sent to user. The following placeholders are supported in this template

* `{{template.otp}}` - Generated OTP.
* `{{template.otpExpiryPeriod}}` - How long the OTP is valid (in minutes)
* `{{template.phoneNumber}}` - Phone number for verification.

### PASSWORD_FAILED_ATTEMPTS

If there are multiple attempts to change user account password with invalid current password, then the user will get this SMS notification warning user about malicious login attempts. This notification will be sent after 5 invalid attempts. The following placeholders are supported in this template

* `{{user.lockoutPeriod}}` - How long the account will be in locked state (in minutes)

### MFA_OTP

This SMS template is used for login using OTP for multi-fator authentication.
The following placeholders are supported in this template

* `{{template.otp}}` - Generated OTP.
* `{{template.otpExpiryPeriod}}` - How long the OTP is valid (in minutes)

## Example Usage

The following example manages an email template for an org

```hcl
resource "hsdp_iam_sms_template" "mfa_otp_default" {
  type            = "MFA_OTP"
  organization_id = var.org_id

  message = "Hi {{user.givenName}}, your OTP code is {{template.otp}}, valid for {{template.otpExpiryPeriod}} minutes"
}

resource "hsdp_iam_sms_template" "mfa_otp_nl" {
  type            = "MFA_OTP"
  organization_id = var.org_id
  
  message = "Hallo {{user.givenName}}, jouw OTP code is {{template.otp}}, geldig voor {{template.otpExpiryPeriod}} minuten"
  locale  = "NL-nl"
}
```

## Argument Reference

The following arguments are supported:

* `organization_id` - (string, Required) The UUID of the IAM Org to apply this SMS template to
* `type` - (string, Required) The SMS template type. See the `Type` table above for available values
* `message` - (string, Semi-Required) The message, including template placeholders. Max length is 160 chars. Take into account placeholder expansion
* `message_base64` - (base64, Semi-Required) Conflicts with `message`. Same as `message` but provided as a base64 string
* `locale` - (string, Optional) The locale of the template. When not specified the template will become the default. Only a single default template is allowed of course.

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the email template

## Import

Importing is supported but not recommended as the `message` body is not returned when reading out a template via the IAM API
