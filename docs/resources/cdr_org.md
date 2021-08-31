# hsdp_cdr_org

Provides a resource for onboarding HSDP CDR [organizations](https://www.hsdp.io/documentation/clinical-data-repository/stu3/getting-started/ehr).
This resource provides very limited management of the onboarded FHIR organization. At this time off-boarding is also
not supported at API level, so the provider will silently forget a CDR Organization when destroy is called for and will try to
rediscover already onboarded organizations and import their state.

## Example Usage

The following example creates and onboards a CDR FHIR organization

```hcl
data "hsdp_cdr_fhir_store" "sandbox" {
  base_url = "https://cdr-stu3-sandbox.hsdp.io"
  fhir_org_id = var.root_org_id
}

resource "hsdp_cdr_org" "hospital" {
  fhir_store = data.hsdp_cdr_fhir_store.sandbox.endpoint
  org_id = var.sub_org_id

  name = "Hospital"
  part_of = var.root_org_id
  
  purge_delete = false
}
```

## Argument Reference

The following arguments are supported:

* `fhir_store` - (Required) The CDR FHIR store to use
* `org_id` - (Required) The Org ID (GUID) under which to onboard. Usually same as IAM Org ID
* `name` - (Required) The name of the FHIR Org
* `part_of` - (Optional) The parent Organization ID (GUID) this Org is part of
* `purge_delete` - (Optional) If set to `true`, when the resource is destroyed the provider will purge all FHIR resources associated with the Organization. The `ORGANIZATION.PURGE` IAM permission is required for this to work. Default: `false`

!> Only use `purge_delete = true` when you are sure recursive deletion of FHIR resources under the Organization is acceptable for the given deployment.

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the organization

## Import

An existing Organization can be imported using `terraform import hsdp_cdr_org`, e.g.

```bash
terraform import hsdp_cdr_org.myorg a-guid
```
