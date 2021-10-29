---
subcategory: "DICOM Store"
---

# hsdp_dicom_store_config

This resource manages DICOM store configuration of an HSDP provisioned DICOM Store.

## Example usage

The following example demonstrates the basic configuration of a DICOM store

```hcl
resource "hsdp_dicom_store_config" "dicom" {
  config_url = var.dicom_base_url
  organization_id = var.iam_org_id
  
  cdr_service_account {
    service_id = hsdp_iam_service.cdr.service_id
    private_key = hsdp_iam_service.cdr.private_key
  }
  
  fhir_store {
    mpi_endpoint = "https://foo.bar/xxx"      
  }
}
```

## Argument reference

* `config_url` - (Required) The base config URL of the DICOM Store instance
* `organization_id` - (Required) the IAM organization ID to use for authorization
* `cdr_service_account` - (Optional) Details of the CDR service account
  * `service_id` - the service id
  * `private_key` - the service private key
* `fhir_store` - (Optional) the FHIR store configuration
  * `mpi_endpoint` - the FHIR mpi endpoint
  
## Attribute reference

* `qido_url` - QIDO API endpoint URL
* `stow_url` - STOW API endpoint URL
* `wado_url` - WADO API endpoint URL
