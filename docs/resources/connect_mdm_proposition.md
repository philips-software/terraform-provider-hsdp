---
subcategory: "Master Data Management (MDM)"
page_title: "HSDP: hsdp_connect_mdm_proposition"
description: |-
  Manages HSDP Connect MDM Propositions
---

# hsdp_connect_mdm_proposition

Create and manage MDM Proposition resources

~> Currently, deleting Proposition resources is not supported by the MDM API, so use them sparingly

## Example Usage

```hcl
resource "hsdp_connect_mdm_proposition" "app" {
  name        = "moonshot"
  description = "Terraform managed proposition"
  
  organization_id = var.org_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Proposition
* `description` - (Optional) A short description of the Proposition
* `organization_id` - (Required) The ID of the IAM organization this Proposition should fall under
* `status` - (Required) The status of the Proposition [`DRAFT`, `ACTIVE`]

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference (format: `Proposition/${GUID}`)
* `guid` - The GUID of the underlying IAM resource
* `proposition_id` - The ID of the IAM proposition this proposition falls under
* `proposition_guid` - The GUID of the MDM proposition resource
