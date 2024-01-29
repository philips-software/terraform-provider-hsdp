---
subcategory: "Data Broker Service (DBS)"
page_title: "HSDP: hsdp_dbs_subscription_topic"
description: |-
  Manages Connect DBS SQS Subscription topics
---

# hsdp_dbs_topic_subscription

Manages Connect DBS Topic Subscriptions

## Example Usage

```hcl
resources "hsdp_dbs_topic_subscription" "my-topic" {
  name_infix  = "my-topic"
  description = "My topic"
  
  data_type     = "some-type"
  subscriber_id = hsdp_dbs_sqs_subscriber.my-subscriber.id
  
  deliver_data_only = true
  
}
```

## Argument Reference

The following arguments are supported:

* `name_infix` - (Required) The name infix of the subscription
* `description` - (Required) A short description of the subscription
* `subscriber_id` - (Required) The ID of the subscriber to associate with this topic subscription
* `data_type` - (Required) The data type of the topic
* `deliver_data_only` - (Optional) Boolean designating whether to deliver only data (true) or data and metadata (false). Default is `false`
* `kinesis_stream_partition_key` - (Optional) When used in combination with a Kinesis Subscriber, the Stream Partition Key for inserting the data into the Kinesis Stream needs to be provided. Example:= `${newuuid()}`

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the subscription (format: `${GUID}`)
* `name` - The name of the subscription
* `status` - The status of the subscription
* `rule_name` - The rule name of the subscription

## Import

DBS Topic Subscriptions cannot be imported into this resource.
