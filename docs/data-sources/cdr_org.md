---
subcategory: "Clinical Data Repository (CDR)"
---

# hsdp_cdr_org

Retrieve details of a CDR organization

## Example Usage

The following example retrieves details of an onboarded CDR org

```hcl
data "hsdp_cdr_fhir_store" "sandbox" {
  base_url = "https://cdr-stu3-sandbox.hsdp.io"
  fhir_org_id = var.root_org_id
}

data "hsdp_cdr_org" "hospital" {
  fhir_store = data.hsdp_cdr_fhir_store.sandbox.endpoint
  org_id = var.sub_org_id

  # Set up this org to use FHIR R4
  version = "r4"
}

output "cdr_org_name" {
  value = data.hsdp_cdr_org.hospital.name
}
```

## Argument Reference

The following arguments are supported:

* `fhir_store` - (Required) The CDR FHIR store to use
* `version` - (Optional) The FHIR version to use. Options [ `stu3` | `r4` ]. Default is `stu3`
* `org_id` - (Required) The Org ID (GUID) under which to onboard. Usually same as IAM Org ID

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the organization
* `name` - The name of the FHIR Org
* `part_of` - The parent Organization ID (GUID) this Org is part of
