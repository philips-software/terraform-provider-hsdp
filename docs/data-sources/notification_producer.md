---
subcategory: "Notification"
---

# hsdp_notification_producer

Look up a Notification producer

## Example usage

```hcl
data "hsdp_notification_producer" "producer" {
  producer_id =  "example-d8f5-4fe4-b486-29a7fd30c9ba"
}
```

## Argument reference

* `producer_id` - (Required) The UUID of the IAM producer

## Attribute reference

* `managing_organization` - The name of IAM organization or tenant
* `producer_product_name` -  TThe name of the product
* `producer_service_name` - The name of the service within the product
* `producer_service_instance_name` - The name of a service instance in the product.
* `producer_service_base_url` - The base URL of the producer
* `producer_service_path_url` - The URL extension of the producer
* `description` - Description of the producer application
