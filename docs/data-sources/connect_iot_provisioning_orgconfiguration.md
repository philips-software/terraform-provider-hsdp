---
subcategory: "Connect IoT"
---

# Data Source: hsdp_connect_iot_provisioning_orgconfiguration

Retrieves an existing Connect IoT Provisioning organization configuration.

## Example Usage

```hcl
data "hsdp_connect_iot_provisioning_orgconfiguration" "myconfig" {
  organization_guid = var.organization_guid
}

output "public_key" {
  value = data.hsdp_connect_iot_provisioning_orgconfiguration.myconfig.public_key
}
```

## Argument Reference

The following arguments are required:

* `organization_guid` - (Required) The organization GUID to look up.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the organization configuration.
* `service_account` - Service account configuration block containing:
  * `client_id` - The service account client ID.
  * `service_account_key` - (Sensitive) The service account key.
  * `token_url` - The token URL for the service account.
* `bootstrap_signature` - Bootstrap signature configuration block containing:
  * `algorithm` - The signing algorithm used.
  * `public_key` - The public key for signature verification.
  * `salt_length` - The salt length for the signature.
