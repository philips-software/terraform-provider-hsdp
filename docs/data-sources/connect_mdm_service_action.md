---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_service_action

Retrieve details of a ServiceAction

## Example Usage

```hcl
data "hsdp_connect_mdm_service_action" "delete_bucket_action" {
   name = "delete-bucket"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service action

## Attributes Reference

The following attributes are exported:

* `id` - The ServiceAction ID
* `description` - The ServiceAction description
* `type` - The service action type
* `organization_guid` - The service action organization ID
* `standard_service_id` - The standard service ID of this action
* `guid` - The GUID of this service action
