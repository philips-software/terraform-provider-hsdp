---
subcategory: "IAM"
---

# hsdp_iam_application

Provides a resource for managing HSDP IAM application under a proposition.

## Example Usage

The following example creates an application

```hcl
resource "hsdp_iam_application" "testapp" {
  name                = "TESTAPP"
  description         = "Test application"
  proposition_id      = hsdp_iam_proposition.testprop.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the application
* `description` - (Required) The description of the application
* `proposition_id` - (Required) the proposition ID (GUID) to attach this a application to
* `global_reference_id` - (Optional, UUIDv4) Reference identifier defined by the provisioning user. Highly recommend to never set this and let Terraform generate a UUID for you.

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the application

## Import

An existing application can be imported using `terraform import hsdp_iam_application`, e.g.

```shell
terraform import hsdp_iam_application.myapp a-guid
```
