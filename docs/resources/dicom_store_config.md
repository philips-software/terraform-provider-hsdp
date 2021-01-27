# hsdp_dicom_store_config
This resource manages DICOM store configuration of an HSDP provisioned DICOM Store.

# Example usage
The following example demonstrates the basic configuration of a DICOM store

```hcl
resource "hsdp_dicom_store_config" "dicom" {
  base_url = var.dicom_base_url
  
  cdr_service_account {
    service_id = hsdp_iam_service.cdr.service_id
    private_key = hsdp_iam_service.cdr.private_key
  }
  
  fhir_store {
    mpi_endpoint = "https://foo.bar/xxx"      
  }
}
```

