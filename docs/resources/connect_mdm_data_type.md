---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_data_type

Create and manage MDM DataType resources

## Example Usage

```hcl
resource "hsdp_connect_mdm_data_type" "some_type" {
  name        = "tf-some-data-type"
  description = "A Terraform provisioned DataType"
  
  tags = ["ONE", "TWO", "THREE"]
  
  proposition_id = data.hsdp_connect_mdm_proposition.first.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the device group
* `description` - (Optional) A short description of the device group
* `application_id` - (Required) Reference to the Application this group falls under
* `default_iam_group_id` - (Optional) The IAM Group from which this group will inherit roles from

~> The `name` maps to an AWS IoT thing group so this should be globally unique and not used (or re-used) across deployments

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `DataType/${GUID}`)
* `guid` - The GUID of the service action
