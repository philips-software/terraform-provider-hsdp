# hsdp_dicom_store_config
This resource manages DICOM store configuration of an HSDP provisioned DICOM Store.

# Example usage
The following example demonstrates the basic configuration of a DICOM store

```hcl
resource "hsdp_dicom_store_config" "dicom" {
  base_url = var.dicom_base_url
  
  cdr_service_account {
    service_id = "a@b.com"
    private_key = var.service_private_key
  }
  
  fhir_store {
    mpi_endpoint = "https://foo.bar/xxx"      
  }
}
```

