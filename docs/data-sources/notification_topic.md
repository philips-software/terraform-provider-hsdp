---
subcategory: "Notification"
---

# hsdp_notification_topic

Look up a HSDP Notification Topic resource

## Example usage

```hcl
data "hsdp_notification_topic" "topic" {
  name =  "some-topic"
}
```

## Argument reference

* `topic_id` - (Optional) The GUID of the topic to look up.
* `name` - (Optional) The name of the topic look up.

-> Specify either a `topic_id` or a `name`, not both.

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the found topic
* `producer_id` - The UUID of the producer associated with this topic
* `scope` - The scope of this topic. Can be either `public` or `private`
* `allowed_scopes` - The list of allowed scopes
* `is_auditable` -  whether the topic has to be audited whenever messages are published to it.
* `description` - The intended usage of this topic
