---
subcategory: "DICOM Gateway"
page_title: "HSDP: hsdp_dicom_gateway_config"
description: |-
  Manages HSDP DICOM Gateway configurations
---

# hsdp_dicom_gateway_config

This resource manages DICOM Gateway configuration for Store and QueryRetrieve services using HSDP provisioned DICOM Store configuration service.

## Example usage

The following example demonstrates the basic configuration of a DICOM Gateway

```hcl
resource "hsdp_dicom_gateway_config" "dicom_gateway" {
  config_url = var.config_url
  
  organization_id = var.site_id
  
  store_service {
    title = "Store title"
    
    is_secure = false
    port = 104
    
    pdu_length = 65535
    artim_timeout = 3000
    association_idle_timeout = 4500
    
    application_entity {
      allow_any = true
      ae_title = "Foo"
      organization_id = "aaa-bbb-ccc-ddd"
    }
  }

  query_retrieve_service {
    title = "Query Retrieve Title"
    
    is_secure = false
    port = 108
    
    application_entity {
      allow_any = true
      ae_title = "Foo"
      organization_id = "aaa-bbb-ccc-ddd"
    }
  }
}
```

## Argument reference

* `config_url` - (Required) The base config URL of the DICOM Store
* `organization_id` - (Required) The site organization ID
* `store_service` - (Optional) Details of the DICOM Store SCP
  *`title` - Store Device Title
  * `is_secure` - Is secure. Default `false`
  * `port` - Port. Default `104` for Non-Secure and `105` for Secure. Don't change this.
  * `pdu_length` - PDU length. Default `65535`
  * `artim_timeout` - Time-out waiting for A-ASSOCIATE RQ PDU on open TCP/IP connection (Artim timeout). Default `3000 ms`
  * `association_idle_timeout` - Association idle timeout. `4500 ms`
  * `certificate_id` - (Optional) Certificate ID.
    Only needed for secure connections.
  * `authenticate_client_certificate` - (Optional, Boolean) Weather or not the client certificate is authenticated.
    Only needed for secure connections.
  * `application_entity` - Application entity
    * `allow_any` - Allow any. Value can be `true` or `false`
    * `ae_title` - AE title. Allowed characters for aetitle are `A-Za-z0-9\\s/+=_-`. Eg. `DicomStoreScp`
    * `site_organization_id` - Site Organization ID for which Gateway to be deployed

* `query_retrieve_service` - (Optional) the FHIR store configuration
  * `title` - Store Device Title
  * `is_secure` - Is secure. Default `false`
  * `port` - Port. Default `108` for Non-Secure and `109` for Secure. Don't change this.
  * `pdu_length` - PDU length. Default `65535`
  * `artim_timeout` - Time-out waiting for A-ASSOCIATE RQ PDU on open TCP/IP connection (Artim timeout). Default `3000 ms`
  * `association_idle_timeout` - Association idle timeout. `4500 ms`
  * `certificate_id` - (Optional) Certificate ID.
    Only needed for secure connections.
  * `authenticate_client_certificate` - (Optional, Boolean) Weather or not the client certificate is authenticated.
    Only needed for secure connections.
  * `application_entity` - Application entity
    * `allow_any` - Allow any. Value can be `true` or `false`
    * `ae_title` - AE title. Allowed characters for aetitle are `A-Za-z0-9\\s/+=_-`. Eg. `DicomQueryRetrieveScp`
    * `site_organization_id` - Site Organization ID for which Gateway to be deployed
