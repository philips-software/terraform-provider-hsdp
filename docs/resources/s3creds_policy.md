---
subcategory: "S3 Credentials"
page_title: "HSDP: hsdp_s3creds_policy"
description: |-
  Manages HSDP S3 Credentials policies
---

# hsdp_s3creds_policy

-> **Deprecation Notice** This resource is deprecated and will be removed in an upcoming release of the provider

Provides a resource for managing HSDP S3 Credentials policies

> This resource is only available when `credentials_url` is set in the provider config

## Example Usage

The following example creates a new policy

```hcl
resource "hsdp_s3creds_policy" "policy1" {
  product_key = var.credentials_product_key

  policy = <<POLICY
{
  "conditions": {
    "managingOrganizations": [ "${var.org_id}" ],
    "groups": [ "PublishGroup" ]
  },
  "allowed": {
    "resources": [ "${var.org_id}/foo/*" ],
    "actions": [
      "GET",
      "PUT",
      "LIST",
      "DELETE"
    ]
  }
}
POLICY
}
```

## Argument Reference

The following arguments are supported:

* `product_key` - (Required) The product key (tenant) for which this
   policy should apply to
* `policy` - (Required) The policy definition. This is a JSON string as per
   HSDP S3 Credentials documentation

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the policy

## Import

Importing existing policies is currently not supported
