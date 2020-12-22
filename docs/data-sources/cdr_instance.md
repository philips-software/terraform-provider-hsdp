# hsdp_cdr_instance

Retrieve details of an existing Clinical Data Repository (CDR).

## Example Usage

```hcl
data "hsdp_cdr_instance" "mycdr" {
   base_url = "https://sandbox-stu3-cdr.hsdp.io"
}
```

## Argument Reference

The following arguments are supported:

* `base_url` - (Required) the base URL of the CDR instances. This is provided by HSDP.

## Attributes Reference

The following attributes are exported:

* `fhir_store` - The FHIR store base URL
* `type` - The type of CDR deployment. Currently, always `EHR`
