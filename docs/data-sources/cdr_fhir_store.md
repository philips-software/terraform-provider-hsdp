# hsdp_cdr_fhir_store

Retrieve details of an existing Clinical Data Repository (CDR).

## Example Usage

```hcl
data "hsdp_cdr_fhir_store" "mycdr" {
   base_url = "https://sandbox-stu3-cdr.hsdp.io"
   fhir_org_id = var.iam_org_id
}
```

## Argument Reference

The following arguments are supported:

* `base_url` - (Required) the base URL of the CDR instances. This is provided by HSDP.
* `fhir_org_id` - (Required) the FHIR Org ID (GUID)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `endpoint` - The FHIR store endpoint URL
* `type` - The type of CDR deployment. Currently, always `EHR`
