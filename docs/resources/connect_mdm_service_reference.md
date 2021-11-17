---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_service_reference

Create and manage MDM ServiceReference resources

## Example Usage

```hcl
resource "hsdp_connect_mdm_service_reference" "some_reference" {
  name        = "some-service-reference"
  description = "Terraform provisioned service reference"
  
  application_id      = data.hsdp_connect_mdm_application.app.id
  standard_service_id = hsdp_connect_mdm_standard_service.service.id
  matching_rule       = "*"
 
  service_action_ids = [
    hsdp_connect_mdm_service_action.some_action.id
  ]
  
  bootstrap_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service action
* `description` - (Optional) A short description of the service action
* `application_id` - (Required) The application associated with this service reference
* `standard_service_id` - (Required) Reference to a Standard Service
* `matching_rule` - (Required) The rule to use to match up the services
* `service_action_ids` (Required, list(string)) The list of serviced action IDs
* `bootstrap_enabled` (Optional) Wether or not to enable this for bootstrapping

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `ServiceReference/${GUID}`)
* `guid` - The GUID of the service action
