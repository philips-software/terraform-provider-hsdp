---
subcategory: "Identity and Access Management"
---

# hsdp_iam_sms_gateway

Provides a resource for managing HSDP IAM SMS gateway configurations.

## Example Usage

The following example creates an IAM SMS Gateway configuration for an IAM organization

```hcl
resource "hsdp_iam_sms_gateway" "config" {
  organization_id = var.iam_org_id
  
  gateway_provider = "twilio"
  
  properties {
    sid = var.twilio_sub_account_sid
    endpoint = var.twilio_endpoint
    from_number = var.twilio_phone_number
  }
  
  credentials {
    token = var.twilio_sub_account_token
  }
  
  activation_expiry = 15 # OTP is valid for 15 minutes
}
```

## Argument Reference

The following arguments are supported:

* `organization_id` - (Required) the IAM organization ID (GUID) for which this SMS gateway should be
* `gateway_provider` - (Optional) The SMS gateway provider. Default: `twilio`. Supported: [ `twilio` ]
* `properties` - (Required) The properties of the SMS gateway
  * `sid` - (Required) The Twilio sub-account SID (sensitive)
  * `endpoint` - (Required) The Twilio endpoint to use
  * `from_number` - (Required) The Twilio phone number from which SMS messages will appear from
* `credentials` - (Required) Credentials of the SMS gateway
  * `token` - (Required)  The Twilio sub-account token (sensitive)
* `activation_expires_on` - (Optional) Sets the expiry time in minutes of the OTP. Default: `15`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the SMS gateway config

## Import

Existing SMS gateway configurations can be imported.
