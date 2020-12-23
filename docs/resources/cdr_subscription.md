# hsdp_cdr_subscription
Provides a resource for managing [FHIR Subscriptions](https://www.hl7.org/fhir/stu3/subscription.html) in a CDR. 
The only supported channel type is `rest-webhook` therefore the `endpoint` and `headers` are top-level arguments.

## Example Usage

The following example creates a subscription that calls a REST endpoint whenever a Patient resources is changed in CDR

```hcl
data "hsdp_cdr_instance" "sandbox" {
  base_url = "https://cdr-stu3-sandbox.us-east.philips-healthsuite.com"
}

resource "hsdp_cdr_subscription" "patient_changes" {
  fhir_store = data.hsdp_cdr_instance.sandbox.fhir_store
  org_id = var.iam_org_id

  criteria = "Patient"
  endpoint = "https://webhook.myapp.io/patient"
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

* `fhir_store` - (Required) The CDR FHIR store to use
* `org_id` - (Required ) The Org ID of the tenant (GUID) to create the Subscription 
* `criteria` - (Required) On which resource to notify
* `endpoint` - (Required) The REST endpoint to call. Must use `https://`  schema
* `end` - (Required) RFC3339 formatted timestamp when to end notifications
* `reason` - (Optional) Reason for creating the subscription
* `headers` - (Optional) List of headers to add to the REST call

## Attributes Reference

The following attributes are exported:

* `status` - The status of the subscription (requested|active|error|off)

## Import

An existing Subscription can be imported using `terraform import hsdp_cdr_subscription`, e.g.

```bash
terraform import hsdp_cdr_subscription.myorg a-guid
```
