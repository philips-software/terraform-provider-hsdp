# hsdp_iam_password_policy
Provides a resource for managing HSDP IAM [password policies](https://www.hsdp.io/documentation/identity-and-access-management-iam/api-documents#_password_policy). 

## Example Usage

The following example manages a password policy for an org 

```hcl
resource "hsdp_iam_pasword_policy" "mypolicy" {
  managing_organization = data.hsdp_iam_org.myorg.id
  
  expiry_period_in_days = 180
  history_count = 5
  
  complexity {
  	min_length = 8
  	max_length = 32
  	min_numerics = 1
  	min_lowercase = 1
  	min_uppercase = 1
  	min_special_chars = 1  	
  }
}
```

## Argument Reference

The following arguments are supported:

* `managing_organization` - (Required) The UUID of the IAM Org to apply this policy to
* `expiry_period_in_days ` - (Optional) number - The number of days after which the user's password expires.
* `complexity` - (Required) Different criteria that decides the strength of user password for an organization. Block
* `history_count` - (Optional) The number of previous passwords that cannot be used as new password.
* `challenges_enabled` - (Optional) A boolean value indicating if challenges are enabled at organization level. If the value is set to true, `challenge_policy` attribute is mandatory.
* `challenge_policy` - (Mandatory, if `challenges_enabled` = `true`) Specify KBA settings

The `complexity` block supports:

* `min_length` - (Optional) The minimum number of characters password can contain. Default 8
* `max_length` - (Optional) The maximum number of characters password can contain.
* `min_numeric` - (Optional) The minimum number of numerical characters password can contain.
* `min_uppercase` - (Optional) The minimum number of uppercase characters password can contain.
* `min_lowercase` - (Optional) The minimum number of lower characters password can contain.
* `min_special_chars` - (Optional) The minimum number of special characters password can contain.


The `challenge_policy` block supports:

* `default_questions` - (Mandatory) A Multi-valued String attribute that contains one or more default question a user may use when setting their challenge questions.
* `min_question_count` - (Mandatory) An Integer indicating the minimum number of challenge questions a user MUST answer when setting challenge question answers.
* `min_answer_count` - (Mandatory) An Integer indicating the minimum number of challenge answers a user MUST answer when attempting to reset their password.
* `max_incorrect_attempts` - (Mandatory) An Integer indicates the maximum number of failed reset password attempts using challenges.

   
## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the password policy

## Import

If the organization already has a password policy it will be imported automatically.