---
subcategory: "HealthSuite Edge"
---

# hsdp_edge_custom_cert

Manage custom certificates on Edge devices. Set `sync` to true to immediately sync the certificate to the k3s cluster, otherwise
you should create a dependency on a `hsdp_edge_sync` resource to batch sync changes.

## Example usage

```hcl
resource "hsdp_edge_custom_cert" "cert" {
  serial_number = var.serial_number
  
  name = "terrakube.com"
  cert_pem = var.cert_pem
  private_key_pem = var.private_key_pme
}
```

## Argument reference

* `serial_number` - (Required) Device to attach the cert to
* `name` - (Required) Name of the certificate
* `cert_pem`  - (Required) The certificate in PEM format
* `private_key_pem` - (Required) the private key of the certificate in PEM format
* `sync` (Optional, boolean) - When set to true syncs the config after mutations. Default is true.
  Set this to false if you want to batch sync to your device using `hsdp_edge_sync`  
* `principal` - (Optional) The optional principal to use for this resource
  * `uaa_username` - (Optional) The UAA username to use
  * `uaa_password` - (Optional) The UAA password to use
  * `region` - (Optional) Region to use. When not set, the provider config is used
  * `endpoint` - (Optional) The endpoint URL to use if applicable. When not set, the provider config is used

## Attribute reference

* `id` - The id of the custom certificate

## Importing

Importing a custom certificate is supported but not recommended.
