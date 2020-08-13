# hsdp_credentials_policy
Provides a resource for managing HSDP S3 Credentials policies

> This resource is only available when `credentials_url` is set in the provider config

## Example Usage

The following example creates a new policy

```hcl
resource "hsdp_credentials_policy" "policy1" {
  product_key = var.credentials_product_key

  policy = <<POLICY
{
  "conditions": {
    "managingOrganizations": [ var.org_id ],
    "groups": [ "PublishGroup" ]
  },
  "allowed": {
    "resources": [ "${var.org_id}/foo/*" ],
    "actions": [
      "GET",
      "PUT"
    ]
  }
}
POLICY
}
```

## Argument Reference

The following arguments are supported:

* `product_key` - (Required) The product key (tenant) for which this policy should apply to
* `policy` - (Required) The policy definition. This is a JSON string as per HSDP S3 Credentials documentation


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the policy

## Import

Importing existing policies is currently not supported
