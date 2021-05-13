# hsdp_notification_producer
Search for Notification producers

## Example usage

```hcl
data "hsdp_notification_producer" "producer" {
  managing_organization_id =  "example-d8f5-4fe4-b486-29a7fd30c9ba"
  managing_organization =  "exampleOrg"
  producer_product_name =  "exampleProduct"
  producer_service_name = "exampleServiceName"
}
```

## Argument reference
* `managing_organization_id` - (Required) The UUID of the IAM organization or tenant
* `managing_organization` - (Required) The name of IAM organization or tenant
* `producer_product_name` - (Required) TThe name of the product
* `producer_service_name` - (Required) The name of the service within the product

## Attribute reference
* `id` - The producer ID
* `producer_service_instance_name` - The name of a service instance in the product. Used to differentiate multiple copies of the same service used in an organization
* `producer_service_base_url` - The base URL of the producer
* `producer_service_path_url` - The URL extension of the producer
* `description` - Description of the producer application
