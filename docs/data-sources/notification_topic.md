# hsdp_notification_topic
Search for  HSDP Notification Topic resources

## Example usage

```hcl
data "hsdp_notification_topic" "topic" {
  name =  "Topic1"
  producer_id =  "036c8a21-6906-4485-b2e7-e31883d8f9ed"
  scope =  "public"
}
```

## Argument reference
* `name` - (Optional) The name of a topic. The topic name length is restricted to a maximum of 256 characters. The special characters allowed are `-` and `_`.
* `producer_id` - (Optional) The UUID of the producer
* `scope` - (Optional) The scope of the topic to search for

## Attribute reference
* `id` - The topic ID
* `allowed_scopes` - The list of allowed scopes
* `is_auditable` -  whether the topic has to be audited whenever messages are published to it.
* `description` - The intended usage of this topic