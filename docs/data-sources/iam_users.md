---
subcategory: "IAM"
---

# hsdp_iam_users

Search for users based on certain filters

## Example usage

Get all users with unverified email addresses and in disabled state

```hcl
data "resource_hsdp_iam_users" "unactivated" {
  organization_id = var.org_id
  
  email_verified = false
  disabled       = true
}
```

## Argument Reference

The following arguments are supported:

* `organization_id` - (Required) The organization users should belong to
* `email_verified` - (Optional) Filter users on verified email state
* `disabled` - (Optional) Filter users on account disabled status

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ids` - The list of matching users
* `logins` - The list matching user login ids
* `email_addresses` - The email addresses of the matching users
