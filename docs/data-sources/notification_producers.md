# hsdp_notification_producers

Search for Notification producers

## Example usage

```hcl
data "hsdp_notification_producer" "producer" {
  managing_organization_id =  "foo-d8f5-4fe4-b486-29a7fd30c9ba"
}
```

## Argument reference

* `managing_organization_id` - (Required) The UUID of the managing IAM organization for the producers

## Attribute reference

* `managing_organization` - The name of IAM organization or tenant
* `producer_product_name` -  TThe name of the product
* `producer_service_name` - The name of the service within the product
* `producer_service_instance_name` - The name of a service instance in the product.
* `producer_service_base_url` - The base URL of the producer
* `producer_service_path_url` - The URL extension of the producer
* `description` - Description of the producer application
