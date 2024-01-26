---
subcategory: "Data Broker Service (DBS)"
page_title: "HSDP: hsdp_dbs_subscription_topic"
description: |-
  Manages Connect DBS SQS Subscription topics
---

# hsdp_dbs_subscription_topic

Manages Connect DBS SQS Subscription topics

## Example Usage

```hcl
resources "hsdp_dbs_subscription_topic" "my-topic" {
  name  = "my-topic"
  description = "My topic"
  
  subscriber_id = hsdp_dbs_sqs_subscriber.my-subscriber.id
  
  deliver_data_only = true
  
  data_type ="some-type"
  rule_name = "some-rule"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the topic
* `description` - (Required) A short description of the topic
* `subscriber_id` - (Required) The ID of the subscriber to associate with this topic
* `deliver_data_only` - (Optional) Boolean designating whether to deliver only data (true) or data and metadata (false). Default is `false`
* `data_type` - (Required) The data type of the topic
* `rule_name` - (Required) The rule name of the topic

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the topic (format: `Topic/${GUID}`)
* `status` - The status of the topic

## Import

DBS SQS Subscription topics cannot be imported into this resource.
