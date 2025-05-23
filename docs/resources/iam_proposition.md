---
subcategory: "Identity and Access Management (IAM)"
page_title: "HSDP: hsdp_iam_proposition"
description: |-
  Manages HSDP IAM Proposition resources
---

# hsdp_iam_proposition

Provides a resource for managing HSDP IAM proposition belonging to an Organization.

## Example Usage

The following example creates an application

```hcl
resource "hsdp_iam_proposition" "testprop" {
  name                = "TestProposition"
  description         = "Test Proposition"
  organization_id     = hsdp_iam_org.devorg.id
  wait_for_delete     = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the application
* `description` - (Required) The description of the application
* `organization_id` - (Required) the organization ID (GUID) to attach this a proposition to
* `global_reference_id` - (Optional, UUIDv4) Reference identifier defined by the provisioning user. Highly recommend to never set this and let Terraform generate a UUID for you.
* `wait_for_delete` - (Optional, boolean) If set to true, the resource will wait for the proposition to be deleted before continuing. Default is true.

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the proposition

## Import

An existing proposition can be imported using `terraform import hsdp_iam_proposition`, e.g.

```shell
terraform import hsdp_iam_proposition.myprop a-guid
```
