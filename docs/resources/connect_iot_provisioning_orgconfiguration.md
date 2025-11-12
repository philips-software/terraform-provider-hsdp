---
subcategory: "Provisioning"
page_title: "HSDP: hsdp_connect_iot_provisioning_orgconfiguration"
description: |-
  Manages HSDP DIP Connect IoT Provisioning Org Configuration
---
# hsdp_connect_iot_provisioning_orgconfiguration (Resource)

Provides a resource for managing HSDP Connect IoT provisioning organization configurations. This resource allows you to configure the service account and bootstrap signature settings for an organization in the Connect IoT provisioning service.

## Example Usage

```terraform
resource "hsdp_connect_iot_provisioning_orgconfiguration" "my-orgconfig" {
  organization_guid = "1ac2e233-8146-4661-8ec4-dc956aeb5a4b"
  
  service_account {
    service_account_id  = "demo_test_tf.xyz-app.xyz-prop@demo.iot__connect__sandbox.apmplatform.philips-healthsuite.com"
    service_account_key = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowBlahBlahI1KExUm\n-----END RSA PRIVATE KEY-----"
  }

  bootstrap_signature {
    algorithm  = "RSA-SHA256"
    public_key = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSBlahBlahBlahMqgzQIDAQAB\n-----END PUBLIC KEY-----"
    
    config {
      type        = "RSA"
      padding     = "RSA_PKCS1_PSS_PADDING"
      salt_length = "RSA_PSS_SALTLEN_MAX_SIGN"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `organization_guid` - (Required, ForceNew) The GUID of the organization to configure.

* `service_account` - (Required) Service account configuration block containing:
  * `service_account_id` - (Required) The service account ID for the organization.
  * `service_account_key` - (Required, Sensitive) The service account private key.

* `bootstrap_signature` - (Required) Bootstrap signature configuration block containing:
  * `algorithm` - (Required) The signature algorithm to use (e.g., "RSA-SHA256").
  * `public_key` - (Required) The public key for bootstrap signature verification.
  * `config` - (Optional) Additional configuration block containing:
    * `type` - (Optional) The signature type (e.g., "RSA", "ECC", "DSA").
    * `padding` - (Optional) The padding type (e.g., "RSA_PKCS1_PSS_PADDING").
    * `salt_length` - (Optional) The salt length configuration (e.g., "RSA_PSS_SALTLEN_MAX_SIGN").

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the organization configuration.

## Import

Organization configurations can be imported using their ID:

```shell
terraform import hsdp_connect_iot_provisioning_orgconfiguration.my-orgconfig <id>
```
