# hsdp_dicom_gateway_config
This resource manages DICOM Gateway configuration for Store and QueryRetrieve services using HSDP provisioned DICOM Store configuration service.

# Example usage
The following example demonstrates the basic configuration of a DICOM Gateway

```hcl
resource "hsdp_dicom_gateway_config" "dicom_gateway" {
  config_url = var.config_url
  
  store_service {
    is_secure = false
    port = 104
    
    pdu_length = 65535
    artim_timeout = 3000
    association_idle_timeout = 4500
    
    application_entity {
      allow_any = true
      ae_title = "Foo"
      organization_id = "aaa-bbb-ccc-ddd"
      service_timeout = 0
    }

    application_entity {
      allow_any = true
      ae_title = "Bar"
      organization_id = "bbb-ccc-ddd-eee-bbb"
      service_timeout = 0
    }
  }

  query_retrieve_service {
    is_secure = false
    port = 108
    
    application_entity {
      allow_any = true
      ae_title = "Foo"
      organization_id = "aaa-bbb-ccc-ddd"
      service_timeout = 0
    }

    application_entity {
      allow_any = true
      ae_title = "Bar"
      organization_id = "bbb-ccc-ddd-eee-bbb"
      service_timeout = 0
    }
  }
}
```

# Argument reference

* `config_url` - (Required) The base config URL of the DICOM Store
* `organization_id` - (Required) The site organization ID
* `title` - Store Device Title
* `store_service` - (Optional) Details of the DICOM Store SCP
  * `is_secure` - Is secure. Default `false`
  * `port` - Port. Default `104` for Non-Secure and `105` for Secure. Don't change this.
  * `pdu_length` - PDU length. Default `65535`
  * `artim_timeout` - Time-out waiting for A-ASSOCIATE RQ PDU on open TCP/IP connection (Artim timeout). Default `3000 ms`
  * `association_idle_timeout` - Association idle timeout. `4500 ms`
  * `application_entity` - Application entity
    * `allow_any` - Allow any. Value can be `true` or `false`
    * `ae_title` - AE title. Allowed characters for aetitle are `A-Za-z0-9\\s/+=_-`. Eg. `DicomStoreScp`
    * `site_organization_id` - Site Organization ID for which Gateway to be deployed

* `queryretrieve_service` - (Optional) the FHIR store configuration
  * `is_secure` - Is secure. Default `false`
  * `port` - Port. Default `108` for Non-Secure and `109` for Secure. Don't change this.
  * `pdu_length` - PDU length. Default `65535`
  * `artim_timeout` - Time-out waiting for A-ASSOCIATE RQ PDU on open TCP/IP connection (Artim timeout). Default `3000 ms`
  * `association_idle_timeout` - Association idle timeout. `4500 ms`
  * `application_entity` - Application entity
    * `allow_any` - Allow any. Value can be `true` or `false`
    * `ae_title` - AE title. Allowed characters for aetitle are `A-Za-z0-9\\s/+=_-`. Eg. `DicomQueryRetrieveScp`
    * `site_organization_id` - Site Organization ID for which Gateway to be deployed

# Attribute reference
