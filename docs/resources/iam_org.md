# hsdp_iam_org
Provides a resource for managing HSDP IAM [organizations](https://www.hsdp.io/documentation/identity-and-access-management-iam/concepts/iam-resource-model). 

## Example Usage

The following example creates an org 

```hcl
resource "hsdp_iam_org" "testorg" {
  name          = "TestOrg"
  description   = "Test Organization"
  parent_org_id = hsdp_iam_org.myorg.id
}
```

Assuming the following Org exists or has been imported

```hcl
resource "hsdp_iam_org" "myorg" {
  name        = "MyOrg"
  description = "My IAM Organization"
  is_root_org = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Org in IAM
* `description` - (Required) The description of the Org
* `parent_org_id` - (Required if not root org) The parent Org ID (GUID)
* `is_root_org` - (Optional) Marks the Org as a root organization (boolean)

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the organization

## Import

An existing Organization can be imported using `terraform import hsdp_iam_org`, e.g.

```bash
terraform import hsdp_iam_org.myorg guid4-of-the-org-you-want-to-import-here
```
