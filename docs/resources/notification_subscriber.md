# hsdp_notification_subscriber
Create and manage HSDP Notification subscriber resources

## Example usage

```hcl
resource "hsdp_notification_subscriber" "subscriber" {
  managing_organization_id =  "example-d8f5-4fe4-b486-29a7fd30c9ba"
  managing_organization =  "exampleOrg"
  subscriber_product_name = "subsciberProd"
  subscriber_service_name = "subsciberService"
  subscriber_service_instance_name = "serviceInsttest12"
  subscriber_service_base_url = "https://ns-client-logdev.cloud.pcftest.com/"
  subscriber_service_path_url = "core"
  description = "subscriber description"
}
```

## Argument reference
* `managing_organization_id` - (Required) The UUID of the IAM organization or tenant
* `managing_organization` - (Required) The name of IAM organization or tenant
* `subscriber_product_name` - (Required) The name of the product
* `subscriber_service_name` - (Required) The name of the subscriber service
* `subscriber_service_instance_name` - (Required) The name of a service instance, used to differentiate multiple copies of the same service used in an organization
* `subscriber_service_base_url` - (Required) The base URL of the subscriber
* `subscriber_service_path_url` - (Required) The URL extension of the subscriber
* `description` - (Optional) Description of the subscriber application

## Attribute reference
* `id` - The subscriber ID

## Importing
Importing a HSDP Notification producer is supported.