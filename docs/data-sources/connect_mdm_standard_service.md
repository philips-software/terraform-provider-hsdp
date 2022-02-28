---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_standard_service

Retrieve details of a StandardService

-> The `MDM-STANDARDSERVICE.READ` IAM permissions is required to use this data source

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

## Argument Reference

The following arguments are supports:

* `name` - (Required) The name standard service to look up

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The StandardService ID
* `description` - The StandardService descriptions
* `trusted` - If the service is a trust one
* `service_urls` - The list of URLs
* `organization_id` - To which ORG this service is associated to
* `tag` - The tag associated with this StandardService
