# hsdp_cdr_org
Provides a resource for onboarding HSDP CDR [organizations](https://www.hsdp.io/documentation/clinical-data-repository/stu3/getting-started/ehr).
This resource provides very limited management of the onboarded FHIR organization. At this time offboarding is also
not supported at API level, so the provider will silently forget a CDR Organization when a destroy is called for and will try to 
rediscover already onboarded organizations and import their state.

## Example Usage

The following example creates and onboards a CDR FHIR organization

```hcl
data "hsdp_cdr_fhir_store" "sandbox" {
  base_url = "https://cdr-stu3-sandbox.us-east.philips-healthsuite.com"
  root_org_id = var.iam_org_id
}

resource "hsdp_cdr_org" "hospital" {
  fhir_store = data.hsdp_cdr_fhir_store.sandbox.endpoint
  org_id = var.iam_org_id

  name = "Hospital"
}
```

## Argument Reference

The following arguments are supported:

* `fhir_store` - (Required) The CDR FHIR store to use
* `org_id` - (Required) The Org ID (GUID) under which to onboard. Usually same as IAM Org ID
* `name` - (Required) The name of the FHIR Org
* `part_of` - (Optional) The parent Organization ID (GUID) this Org is part of

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the organization

## Import

An existing Organization can be imported using `terraform import hsdp_cdr_org`, e.g.

```bash
terraform import hsdp_cdr_org.myorg a-guid
```
