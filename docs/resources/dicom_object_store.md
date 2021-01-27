# hsdp_dicom_object_store
This resource manages DICOM Object stores

# Example Usage

```hcl
resource "hsdp_dicom_object_store" "store1" {
  base_url = var.dicom_base_url
  
  description = "Store 1"
  
  static_access {
    endpoint = "https://s3-external.amazonaws.com"
    bucket_name = "xxxx-xxxx-xxxx-xxxx"
    access_key = "xxx"
    secret_key = "yyy"
  }
  
  s3creds_access {
    endpoint = "https://xxx.com"
    product_key = "xxxx-xxxx-xxxx-xxxx"
    bucket_name = "yyyy-yyyy-yyy-yyyy"
    folder_path = "/store1"
    service_account {
      service_id = "a@b.com"
      private_key = var.service_private_key
    }
  }
}
```
