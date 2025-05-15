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
  expiration   = "2025-12-31T23:59:59Z"  # Key valid until the end of 2025
  salt         = "static-salt-value"      # Optional salt for deterministic key generation
  region       = "us-east"
  environment  = "prod"
}

# Examples of other valid time formats
resource "hsdp_tenant_key" "short_lived" {
  project      = "my-project"
  organization = "my-organization"
  signing_key  = "temporary-key"
  expiration   = "2023-06-30T12:00:00Z"  # Mid-year expiration
  salt         = "salt-1"
}

resource "hsdp_tenant_key" "long_lived" {
  project      = "my-project"
  organization = "my-organization"
  signing_key  = "permanent-key"
  expiration   = "2030-01-01T00:00:00Z"  # Valid until 2030
  salt         = "salt-2" 
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The project identifier the key is associated with.
* `organization` - (Required) The organization identifier the key is associated with.
* `signing_key` - (Required) The signing key used to generate the API key. This is sensitive information.
* `scopes` - (Optional) A list of scopes to associate with the key. These define the permissions granted to the key.
* `expiration` - (Optional) The expiration time for the key in RFC3339 format (e.g., "2025-12-31T23:59:59Z"). Default is one year from the current date.
* `salt` - (Optional) A salt value used to generate deterministic API keys. Using the same salt value with the same inputs will always generate the same key.
* `region` - (Optional) The region for the key. Default is "us-east".
* `environment` - (Optional) The environment for the key. Default is "prod".

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The signature of the generated API key, used to identify this resource.
* `signature` - The signature of the generated API key (same as id, for consistency with other resources).
* `key` - The generated API key (sensitive value).

## Import

Tenant keys can be imported using the signature, e.g.,

```shell
terraform import hsdp_tenant_key.example [signature]
```

## Notes

* All fields are marked as ForceNew, meaning any changes to these fields will create a new resource (and destroy the old one).
* The API key is generated deterministically when a `salt` value is provided. Using the same inputs (project, organization, salt, etc.) will always generate the same key.
* The API key generation is no longer time-dependent as it uses a specific expiration time rather than a relative duration.
* Both the API key and signature are computed values based on your configuration.
* The `signature` attribute is provided both as the resource ID and as a separate output field, making it easier to reference in other resources or outputs.