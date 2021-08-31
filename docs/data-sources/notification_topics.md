# hsdp_notification_topics

Search for  HSDP Notification Topic resources

## Example usage

```hcl
data "hsdp_notification_topics" "topic1" {
  name =  "Topic1"
}
```

## Argument reference

* `name` - (Required) The name of a topic. The topic name length is restricted to a maximum of 256 characters. The special characters allowed are `-` and `_`

## Attribute reference

* `topic_ids` - The list of matching topic IDs
