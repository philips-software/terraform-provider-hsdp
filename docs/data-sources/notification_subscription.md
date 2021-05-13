# hsdp_notification_subscription
Looks up  HSDP Notification subscription resources

## Example usage

```hcl
data "hsdp_notification_subscription" "subscription" {
  topic_id = "ca1e3aa4-1409-4b1b-95e5-8795ecdecea7"
  subscriber_id = "4e2546a3-b162-47d1-8014-c89148def43f"
  subscription_endpoint = "https://ns-client-logdev.cloud.pcftest.com/core/notification/Message"
}
```

## Argument reference
* `id` = (Optional) The UUID of the subscription
* `managed_organization_id` - (Optional) The managed organization id
* `managed_organization` - (Optional) The managed organization name
* `topic_id` - (Optional) The UUID of the topic
* `subscriber_id` - (Optional) The UUID of the subscriber

## Attribute reference
* `id` - The subscription ID
* `subscription_endpoint` - The subscription endpoint