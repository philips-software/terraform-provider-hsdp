---
subcategory: "Identity and Access Management (IAM)"
---

# hsdp_iam_permission

Retrieves information on a permission

## Example Usage

```hcl
data "hsdp_iam_permission" "patient_purge" {
  name = "PATIENT.PURGE"
}
```

```hcl
output "patient_purge_description" {
   value = data.hsdp_iam_permission.patient_purge.description
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The ID the permission
* `description` - The description
* `type` - The type of the permission
* `category` - The category this permission is in
