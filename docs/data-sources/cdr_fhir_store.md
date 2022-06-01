---
subcategory: "Clinical Data Repository (CDR)"
---

# hsdp_cdr_fhir_store

Retrieve details of an existing Clinical Data Repository (CDR)

## Example Usage

```hcl
data "hsdp_cdr_fhir_store" "mycdr" {
   base_url = "https://sandbox-stu3-cdr.hsdp.io/store/fhir"
   fhir_org_id = var.iam_org_id
}
```

## Argument Reference

The following arguments are supported:

* `base_url` - (Required) the base URL of the CDR instances. This is provided by HSDP.
* `fhir_org_id` - (Required) the FHIR Org ID (GUID)

~> Earlier versions of this data source required the `base_url` to be specified without the `/store/fhir` path component.
   This is now mandatory, but the data source will append this internally and emit a warning for now if it is missing.
   When upgrading to this or newer versions of the provider please add the path component.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `endpoint` - The FHIR store endpoint URL
* `type` - The type of CDR deployment. Currently, always `EHR`
