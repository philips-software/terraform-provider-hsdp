---
subcategory: "Clinical Data Repository (CDR)"
page_title: "HSDP: hsdp_cdr_subscription"
description: |-
  Manages HSDP CDR Subscription resources
---

# hsdp_cdr_subscription

Provides a resource for managing [FHIR Subscriptions](https://www.hl7.org/fhir/stu3/subscription.html) in a CDR.
The only supported channel type is `rest-webhook` therefore the `endpoint` and `headers` are top-level arguments.

## Example Usage

The following example creates a subscription that calls a REST endpoint whenever a Patient resources is changed in the CDR FHIR store

```hcl
data "hsdp_cdr_fhir_store" "sandbox" {
  base_url = "https://cdr-stu3-sandbox.hsdp.io"
  fhir_org_id = var.org_id
}

resource "hsdp_cdr_org" "test" {
  fhir_store  = data.hsdp_cdr_fhir_store.sandbox.endpoint
  version     = "r4"
  
  name        = "Hospital"
  org_id      = hsdp_iam_org.test.id
}

resource "hsdp_cdr_subscription" "patient_changes" {
  # Refer to the fhir_store through the CDR Organization
  fhir_store = hsdp_cdr_org.test.fhir_store
  version    = "r4"
  
  criteria        = "Patient"
  reason          = "Notification for patient changes"
  endpoint        = "https://webhook.myapp.io/patient"
  delete_endpoint = "https://webhook.myapp.io/patient_deleted"
  headers = [
    "Authorization: Basic cm9uOnN3YW5zb24="
  ]

  end = "2030-12-31T23:59:59Z"
}
```

CDR will send a `POST` request to the endpoint with a JSON body containing:

```json
{
    "logicalId": "df08e38a-4ac7-4434-bca9-479aaab32585",
    "versionId": "df08e38a-4ac7-4434-bca9-479aaab32585",
    "resourceType": "Patient"
}
```

## Argument Reference

The following arguments are supported:

* `fhir_store` - (Required) The CDR FHIR store endpoint to use

~> It is highly recommended to refer to the `fhir_store` attribute of the CDR Organization.
   This creates an explicit dependency between the subscription and the FHIR organization,
   ensuring proper lifecycle handling by Terraform

* `version` - (Optional) The FHIR version to use. Options [ `stu3` | `r4` ]. Default is `stu3`
* `criteria` - (Required) On which resource to notify
* `reason` - (Required) Reason for creating the subscription
* `endpoint` - (Required) The REST endpoint to call. Must use `https://`  schema
* `end` - (Required) RFC3339 formatted timestamp when to end notifications
* `delete_endpoint` - (Optional) The REST endpoint to call for DELETE operations. Must use `https://` schema  
* `headers` - (Optional) List of headers to add to the REST call

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the CDR subscription
* `status` - The status of the subscription (requested | active | error  | off)

## Import

An existing Subscription can be imported using `terraform import hsdp_cdr_subscription`, e.g.

```bash
terraform import hsdp_cdr_subscription.myorg a-guid
```
