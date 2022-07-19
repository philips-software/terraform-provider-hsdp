---
subcategory: "Identity and Access Management (IAM)"
---

# hsdp_iam_role_sharing_policy

Provides a resource for managing HSDP IAM Role Sharing Policies, introduced in the
March 2022 release.

A principal (user / identity) with any of the following permissions can create/update the policy:

* `HSDP_IAM_ROLE_SHARE.WRITE`
* `HSDP_IAM_ORGANIZATION.MGMT`

!>  Changing any permissions assigned to a shared role impacts the application behavior across organizations and sometimes may result in application downtime.
Applying a restrictive sharing policy to an organization automatically and recursively removes any existing assignments from all its children - unless the child organization has an overriding policy to retain the assignments. Removal of assignments are permanent and requires re-assignments by the organization administrators

## Example Usage

The following example creates a role sharing policy

```hcl
# Create the role
resource "hsdp_iam_role" "shared" {
  name        = "SOME Role"
  description = "A role we want to share across ORGs"

  permissions = [
    "PATIENT.READ",
    "PRACTITIONER.READ",
  ]

  managing_organization = hsdp_iam_org.my_org.id
}

# Share the role
resource "hsdp_iam_role_sharing_policy" "policy" {
  sharing_policy = "AllowChildren"
  purpose        = "Share SOME role with another organization"
  
  target_organization_id = hsdp_iam_org.another_org.id
}
```

## Argument Reference

The following arguments are supported:

* `sharing_policy` - (Required) The policy to use
  Sharing of a role with a tenant organization can be in one of the following modes:
  * Restricted: - The assignment role to group operation shall check and allow assignment to the groups present in the target organizations. Any assignment operation - both upward and downward organization
      hierarchy - shall fail the API.
  * AllowChildren: - The assignment role to group operations shall check to restrict the assignment to any group in the target or its children organization. Any assignment operation in the parent organization hierarchy shall fail the API.
  * Denied: - The tenant organization cannot make use of this role. Any attempt to assign the role shall fail the API.
* `target_organization_id` - (Required) The target organization UUID to apply this policy for. This can either be a root IAM Org or a subOrg in an existing hierarchy
* `purpose` - (Optional) The purpose of this role sharing policy mapping

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the role sharing policy (also known as `internalId` at the API level)
* `source_organization_id` - The source organization ID
* `role_name` - The role name

## Import

An existing role sharing policy can be imported using `terraform import hsdp_iam_role_sharing_policy`, e.g.

```shell
> terraform import hsdp_iam_role_sharing_policy.mypolicy a-guid
```
