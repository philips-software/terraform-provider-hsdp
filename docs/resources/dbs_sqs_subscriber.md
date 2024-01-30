---
subcategory: "Data Broker Service (DBS)"
page_title: "HSDP: hsdp_dbs_sqs_subscriber"
description: |-
  Manages Connect DBS SQS Subscriber configurations
---

# hsdp_dbs_sqs_subscriber

Manages Connect DBS SQS Subscriber configurations

## Example Usage

```hcl
resource "hsdp_dbs_sqs_subscriber" "my-subscriber" {
  name_infix  = "my-subscriber"
  description = "My subscriber"
  queue_type  = "Standard"
  
  delivery_delay_seconds           = 0
  message_retention_period_seconds = 0
  receive_wait_time_seconds        = 0
  
  server_side_encryption = true
}
```

## Argument Reference

The following arguments are supported:

* `name_infix` - (Required) The name infix of the subscriber
* `description` - (Required) A short description of the subscriber
* `queue_type` - (Required) The type of queue to create [`Standard` | `FIFO`]
* `delivery_delay_seconds` - (Optional) The time in seconds that the delivery of all messages in the queue will be delayed. An integer from 0 to 900 (15 minutes). The default is 0 (zero).
* `message_retention_period_seconds` - (Optional) The number of seconds Amazon SQS retains a message. Integer representing seconds, from 60 (1 minute) to 1209600 (14 days). The default is 345600 (4 days).
* `receive_wait_time_seconds` - (Optional) The time for which a ReceiveMessage call will wait for a message to arrive (long polling) before returning. An integer from 0 to 20 (seconds). The default is 0 (zero).
* `server_side_encryption` - (Optional) Boolean designating whether to enable server-side encryption. Default is `true`

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the subscriber (format: `${GUID}`)
* `name` - The name of the subscriber
* `status` - The status of the subscriber
* `queue_name` - The name of the SQS queue

## Import

An existing SQS Subscriber using `terraform import hsdp_dbs_sqs_subscriber`, e.g.

```bash
terraform import hsdp_dbs_sqs_subscriber.target guid-of-the-subscriber-to-import
```
