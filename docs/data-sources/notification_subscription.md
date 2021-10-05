# hsdp_notification_subscription

Looks up  HSDP Notification subscription resources

## Example usage

```hcl
data "hsdp_notification_subscription" "subscription" {
  subscription_id = "ca1e3aa4-1409-4b1b-95e5-8795ecdecea7"
}
```

## Argument reference

* `id` = (Optional) The UUID of the subscription

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `managed_organization_id` - (Optional) The managed organization id
* `managed_organization` - (Optional) The managed organization name
* `topic_id` - (Optional) The UUID of the topic
* `subscriber_id` - (Optional) The UUID of the subscriber
* `subscription_endpoint` - The subscription endpoint
