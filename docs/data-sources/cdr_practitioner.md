---
subcategory: "Clinical Data Repository (CDR)"
---

# hsdp_cdr_practitioner

Retrieve details of CDR Practitioner resource

## Example Usage

```hcl
data "hsdp_cdr_practitioner" "doc" {
   fhir_store = data.hsdp_cdr_fhir_store.sandbox.endpoint
   
   guid = "abc-def"
}
```

## Argument Reference

The following arguments are supported:

* `fhir_store` - (Required) the base URL of the CDR instance to search in
* `guid` - (Required) the unique FHIR store ID of the Practitioner to look up

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `fhir_json` - The full FHIR JSON record as returned by CDR
