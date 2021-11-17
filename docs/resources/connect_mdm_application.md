---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_application

Create and manage MDM Application resources

~> Currently, deleting Application resources is not supported by the MDM API, so use them sparingly

## Example Usage

```hcl
resource "hsdp_connect_mdm_application" "app" {
  name        = "mobile"
  description = "Terraform managed Application"
  
  organization_id = var.org_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Application
* `description` - (Optional) A short description of the Application
* `proposition_id` - (Required) The ID of the Proposition this Application should fall under
* `global_reference_id` - (Optional) A global reference ID for this application

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference (format: `Application/${GUID}`)
* `guid` - The GUID of this resource
