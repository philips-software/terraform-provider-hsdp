---
subcategory: "IAM"
---

# hsdp_iam_activation_email

Re-sends activation emails. This resource can be used in combination with the
`hsdp_iam_users` data source to automatically resend activation emails

~> This resource only works when `HSDP_SHARED_KEY` and `HSDP_SHARED_SECRET` are configured or equivalent provider attributes are set. The relevant API requires `API signing`.

## Example usage

```hcl
// Fetch unverified users using remote state
// This trick helps manage the dependencies, even for local
data "terraform_remote_state" "local" {
  backend = "local"
  
  defaults = {
    unverified_users = []
  }
}

// Resend activation email every 168 hours (7 days)
resource "hsdp_iam_activation_email" "activations" {
  for_each = toset(data.terraform_remote_state.local.outputs.unverified_users)
  
  user_id      = for_each.key
  resend_every = 168
}
```

## Argument reference

The following arguments are supported:

* `user_id` - (Required) The user GUID of the user
* `resend_every` - (Optional) Resend activation email after the provided number of hours. Default `72` (3 days)

> Email are only triggered when Terraform is run. When using this resource
> it's best to schedule Terraform to run every resend interval

## Attributes Reference

The following attributes are exported:

* `last_sent` - (Computed) When the last email was sent
* `send` - (Computed) Enabled when activation email is going to be sent
