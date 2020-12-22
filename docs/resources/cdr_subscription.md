# hsdp_cdr_subscription
Provides a resource for subscribing to HSDP CDR [subscriptions](https://www.hsdp.io/documentation/clinical-data-repository/stu3/getting-started/ehr).

## Example Usage

The following example creates a subscription that calls a REST endpoint whenever a Patient resources is changed in CDR

```hcl
data "hsdp_cdr_instance" "sandbox" {
  base_url = "https://cdr-stu3-sandbox.us-east.philips-healthsuite.com"
}

resource "hsdp_cdr_subscription" "patient_changes" {
  fhir_store = data.hsdp_cdr_instance.sandbox.fhir_store
  root_org_id = var.iam_org_id

  criteria = "Patient"
  endpoint = "https://webhook.myapp.io/patient"
  headers = [
    "Authorization: Basic cm9uOnN3YW5zb24="
  ]
  
  end = "2030-12-31T23:59:59Z"
}
```

The REST endpoint will be called with a JSON body as follows:

```json
{
    "logicalId": "df08e38a-4ac7-4434-bca9-479aaab32585",
    "versionId": "df08e38a-4ac7-4434-bca9-479aaab32585",
    "resourceType": "Patient"
}
```

## Argument Reference

The following arguments are supported:

* `criteria` - (Required) On which resource to notify
* `end` - (Required) RFC3339 formatted timestamp when to end notifications
* `reason` - (Optional) Reason for the notification
* `endpoint` - (Required) The REST endpoint to call
* `headers` - (Optional) List of headers to add to call
* `fhir_store` - (Required) The CDR FHIR store to use
* `root_org_id` - (Required ) The root Org ID (GUID) to onboard the organization under

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the subscription

## Import

An existing Subscription can be imported using `terraform import hsdp_cdr_subscription`, e.g.

```bash
terraform import hsdp_cdr_subscription.myorg a-guid
```
