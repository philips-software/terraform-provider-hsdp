---
subcategory: "Master Data Management (MDM)"
page_title: "HSDP: hsdp_connect_mdm_blobl_subscription"
description: |-
  Manages HSDP Connect MDM Blob subscriptions
---

# hsdp_connect_mdm_blob_subscription

Create and manage MDM BlobSubscription resources

## Example Usage

```hcl
resource "hsdp_connect_mdm_blob_subscription" "first" {
  name        = "tf-blob-subscription"
  description = "Terraform provisioned Blob subscription"

  data_type_id     = hsdp_connect_mdm_data_type.first.id
  
  notification_topic_id = var.topic_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Blob subscription
* `description` - (Optional)
* `data_type_id` - (Required)
* `notification_topic_id` - (Required)

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `Group/${GUID}`)
* `guid` - The GUID of the service action
