# hsdp_iam_proposition

Retrieve details of an existing proposition

## Example Usage

```hcl
data "hsdp_iam_proposition" "my_prop" {
   name = "MYPROPOSITION"
   organization_id = var.my_org_id
}
```

```hcl
output "my_prop_display_name" {
   value = data.hsdp_iam_org.display_name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the proposition to look up
* `organization_id` - (Required) the UUID of the organization the proposition belongs to

## Attributes Reference

The following attributes are exported:

* `description` - The description of the organization
* `global_reference_id` - The global reference ID of the proposition
