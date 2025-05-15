---
subcategory: "Tenant Management"
page_title: "HSDP: hsdp_tenant_key"
description: |-
  Manages HSDP Tenant Keys
---

# hsdp_tenant_key

The `hsdp_tenant_key` resource provides a way to generate and manage API keys for accessing HSDP tenant services. These keys are generated with specific scopes, expiration times, and are tied to a specific project and organization.

## Example Usage

```hcl
resource "hsdp_tenant_key" "example" {
  project      = "my-project"
  organization = "my-organization"
  signing_key  = "your-secure-signing-key"
  scopes       = ["scope1", "scope2"]
  expiration   = "24h"  # Key valid for 24 hours
  region       = "us-east"
  environment  = "prod"
}

# Examples of other valid duration formats
resource "hsdp_tenant_key" "short_lived" {
  project      = "my-project"
  organization = "my-organization"
  signing_key  = "temporary-key"
  expiration   = "30m"  # 30 minutes
}

resource "hsdp_tenant_key" "one_week" {
  project      = "my-project"
  organization = "my-organization"
  signing_key  = "weekly-key"
  expiration   = "168h"  # 7 days (1 week)
}

resource "hsdp_tenant_key" "complex_duration" {
  project      = "my-project"
  organization = "my-organization"
  signing_key  = "complex-key"
  expiration   = "72h30m"  # 72 hours and 30 minutes
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The project identifier the key is associated with.
* `organization` - (Required) The organization identifier the key is associated with.
* `signing_key` - (Required) The signing key used to generate the API key. This is sensitive information.
* `scopes` - (Optional) A list of scopes to associate with the key. These define the permissions granted to the key.
* `expiration` - (Optional) The expiration time for the key in Go duration format. Default is "8760h" (approximately 1 year). Examples of valid values include "30m", "24h", "168h" (1 week), "8760h" (1 year). Negative durations are not allowed. Go duration format supports combinations of hours (h), minutes (m), and seconds (s).
* `region` - (Optional) The region for the key. Default is "us-east".
* `environment` - (Optional) The environment for the key. Default is "prod".

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The signature of the generated API key, used to identify this resource.
* `key` - The signing key that was provided (for convenience in outputs).

## Import

Tenant keys can be imported using the signature, e.g.,

```shell
terraform import hsdp_tenant_key.example [signature]
```

## Notes

* All fields are marked as ForceNew, meaning any changes to these fields will create a new resource (and destroy the old one).