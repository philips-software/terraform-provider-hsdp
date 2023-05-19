---
subcategory: "S3 Credentials"
---

# hsdp_s3creds_policy

Gets information on defined S3 Credential policies

> This resource is only available when `s3creds_url` is set in the provider config

## Example Usage

```hcl
data "hsdp_s3creds_policy" "my_org_policies" {
   product_key = var.product_key
   username = var.iam_login
   password = var.iam_password

   filter {
      managing_org = var.my_org_id
   }
}
```

```hcl
output "s3_credential_policies_my_org" {
   value = data.hsdp_s3creds_policy.my_org_policies.policies
}
```

## Argument Reference

The following arguments are supported:

* `product_key` - (Required) The product key under which to search for policies
* `username` - (Optional) The IAM username to authenticate under
* `password` - (Optional) The password of `username`
* `filter` - (Required) The filter conditions block for selecting policies

### filter options

* `id` - (Optional) The id (uuid) of the policy
* `managing_org` - (Optional) Finds policies under `managing_org` (uuid)
* `group_name` - (Optional) Find policies assigned to this group

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policies` - JSON array of policies found using supplied filter values
