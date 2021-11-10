---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_standard_service

Create and manage MDM Standard Services resources

## Example Usage

```hcl
resource "hsdp_connect_mdm_standard_service" "some_service" {
  name        = "Some service"
  description = "A tenant service that does something"
  
  tags = ["TYCHO"]
  
  service_url {
    url        = "https://tycho-service.hsdp.in"
    sort_order = 1
  }
  
  service_url {
    url        = "https://tycho-service.hsdp.nl"
    sort_order = 2
  }
}
```

## Attributes Reference

The following attributes are exported:

* `name` - (Required) The name of the standard service
* `tags` - (Required, list(string)) A tag associated with the service
* `description` - (Optional) A short description of the service
* `service_url` - (Required) Location of this service. Maximum of `5` items are allowed
  * `url` - (Required) the URL of the service
  * `sort_order` (Required, number) the sorting order
  * `authentication_method_id` - (Optional) The id of the authention method to use

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the standard service (format: `StandardService/${GUID}`)
* `guid` - The GUID of the standard service
