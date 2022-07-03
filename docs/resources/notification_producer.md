---
subcategory: "Notification"
---

# hsdp_notification_producer

Create and manage HSDP Notification producer resources

## Example usage

```hcl
resource "hsdp_notification_producer" "producer" {
  managing_organization_id =  "example-d8f5-4fe4-b486-29a7fd30c9ba"
  managing_organization =  "exampleOrg"
  producer_product_name =  "exampleProduct"
  producer_service_name = "exampleServiceName"
  producer_service_instance_name = "exampleServiceInstance"
  producer_service_base_url = "https://ns-producer.cloud.pcftest.com/"
  producer_service_path_url  = "notification/create"
  description =  "product description"
}
```

## Argument reference

* `managing_organization_id` - (Required) The UUID of the IAM organization or tenant
* `managing_organization_id` - (Required) The name of IAM organization or tenant
* `producer_product_name` - (Required) TThe name of the product
* `producer_service_name` - (Required) The name of the service within the product
* `producer_service_instance_name` - (Required) The name of a service instance in the product. Used to differentiate multiple copies of the same service used in an organization
* `producer_service_base_url` - (Required) The base URL of the producer
* `producer_service_path_url` - (Required) The URL extension of the producer
* `description` - (Optional) Description of the producer application
* `principal` - (Optional) The optional principal to use for this resource
  * `service_id` - (Optional) The IAM service ID
  * `service_private_key` - (Optional) The IAM service private key to use
  * `region` - (Optional) Region to use. When not set, the provider config is used
  * `environment` - (Optional) Environment to use. When not set, the provider config is used
  * `endpoint` - (Optional) The endpoint URL to use if applicable. When not set, the provider config is used

## Attribute reference

* `id` - The producer ID

## Importing

Importing a HSDP Notification producer is supported.
