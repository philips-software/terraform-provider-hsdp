---
subcategory: "IAM"
---

# hsdp_iam_org

Retrieve details of an existing organization

## Example Usage

```hcl
data "hsdp_iam_org" "my_org" {
   organization_id = var.my_org_id
}
```

```hcl
output "my_org_name" {
   value = data.hsdp_iam_org.name
}
```

## Argument Reference

The following arguments are supported:

* `organization_id` - (Required) the UUID of the organization to look up

## Attributes Reference

The following attributes are exported:

* `name` - The name of the organization
* `display_name` - The name of the organization suitable for display.
* `description` - The description of the organization
* `active` - Indicates the administrative status of an organization
* `type` - The organization type e.g. `hospital`
* `parent_id` - Unique ID of the parent organization. If the current organization itself is a domain organization, then the parent value will be returned as `root`.
* `external_id` - External ID defined by client that identifies the organization at client side.
