---
subcategory: "Clinical Data Repository (CDR)"
---

# hsdp_cdr_practitioner

Provides a resource for creating [Practitioner FHIR](https://www.hl7.org/fhir/practitioner.html) resources in CDR.
This resource provides limited management of the Practioner resource.

## Example Usage

The following example creates a FHIR Organisation and a Practitioner associated with it.

```hcl
data "hsdp_cdr_fhir_store" "sandbox" {
  base_url = "https://cdr-stu3-sandbox.hsdp.io"
  fhir_org_id = var.root_org_id
}

resource "hsdp_cdr_org" "hospital" {
  fhir_store = data.hsdp_cdr_fhir_store.sandbox.endpoint
  org_id = var.sub_org_id

  # Set up this org to use FHIR R4
  version = "r4"
  
  name    = "Hospital"
  part_of = var.root_org_id
  
  purge_delete = false
}

resource "hsdp_cdr_practitioner" "practitioner" {
  identifier {
    system = "https://iam.hsdp.io"
    value = "xx-xx"
  }

  name {
    text = "Ron Swanson"
    given = ["Ron"]
    family = "Swanson"
  }
}
```

## Argument Reference

The following arguments are supported:

* `fhir_store` - (Required) The CDR FHIR store endpoint to use

~> It is highly recommended to refer to the `fhir_store` attribute of the CDR Organization.
This creates an explicit dependency between the practitioner and the FHIR organization,
ensuring proper lifecycle handling by Terraform

* `identifier` - (Required) The FHIR identifier block
  * `system` - (Required) The system of the identifier e.g. HSP IAM
  * `value` - (Required) the identifier value e.g. the IAM GUID of the practitioner
  * `use` - (Optional) the use value. Can be `usual`, `temp`, `secondary`, `official`

!> `FHIR` Identifiers might contain PII data which will be stored in the Terraform state.
   Please take this into consideration when using this and other FHIR resources of the provider.

* `name` - (Required) The FHIR name block
  * `text` - (Required) The text representation of the name
  * `given` - (Required, list(string)) The list of given names
  * `first` - (Required) The first name
  * `version` - (Optional) The FHIR version to use. Options [ `stu3` | `r4` ]. Default is `stu3`

!> Switching FHIR versions causes the resource to be replaced, so be careful with this.

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the practitioner
* `version_id` - The version of the resource
* `last_updated` - Last update time

## Import

Importing practitioners is currently not supported
