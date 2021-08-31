# hsdp_notification_topic

Look up a HSDP Notification Topic resource

## Example usage

```hcl
data "hsdp_notification_topic" "topic" {
  topic_id =  "036c8a21-6906-4485-b2e7-e31883d8f9ed"
}
```

## Argument reference

* `topic_id` - (Required) The GUID of the topic to look up.

## Attribute reference

* `producer_id` - The UUID of the producer associated with this topic
* `scope` - The scope of this topic. Can be either `public` or `private`
* `allowed_scopes` - The list of allowed scopes
* `is_auditable` -  whether the topic has to be audited whenever messages are published to it.
* `description` - The intended usage of this topic
