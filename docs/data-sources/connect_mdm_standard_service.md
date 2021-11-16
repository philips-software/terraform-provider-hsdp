---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_standard_service

Retrieve details of a StandardService

## Example Usage

```hcl
data "hsdp_connect_mdm_standard_service" "first" {
  name = "first-service"
}
```

```hcl
output "my_first_service_id" {
   value = data.hsdp_connect_mdm_standard_service.first.id
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The StandardService ID
* `description` - The StandardService descriptions
* `trusted` - If the service is a trust one
* `service_urls` - The list of URLs
* `organization_id` - To which ORG this service is associated to
* `tag` - The tag associated with this StandardService
