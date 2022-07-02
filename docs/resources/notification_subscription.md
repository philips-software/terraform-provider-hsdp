---
subcategory: "Notification"
---

# hsdp_notification_subscription

Create and manage HSDP Notification subscription resources

## Example usage

```hcl
resource "hsdp_notification_subscription" "subscription" {
  topic_id = "ca1e3aa4-1409-4b1b-95e5-8795ecdecea7"
  subscriber_id = "4e2546a3-b162-47d1-8014-c89148def43f"
  subscription_endpoint = "https://ns-client-logdev.cloud.pcftest.com/core/notification/Message"
}
```

## Argument reference

* `topic_id` - (Required) The UUID of the topic
* `subscriber_id` - (Required) The UUID of the subscriber
* `subscription_endpoint` - (Required) The subscription endpoint. Only https protocol is allowed
* `principal` - (Optional) The optional principal to use for this resource
    * `service_id` - (Optional) The IAM service ID
    * `service_private_key` - (Optional) The IAM service private key to use
    * `region` - (Optional) Region to use. When not set, the provider config is used
    * `environment` - (Optional) Environment to use. When not set, the provider config is used
    * `endpoint` - (Optional) The endpoint URL to use if applicable. When not set, the provider config is used

## Attribute reference

*`id` - The subscription ID

## Importing

Importing a HSDP Notification subscription is supported.
