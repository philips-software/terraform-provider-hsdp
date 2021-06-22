# hsdp_dicom_gateway_config
This resource manages DICOM gateway configuration of an HSDP provisioned DICOM Store.

# Example usage
The following example demonstrates the basic configuration of a DICOM gateway

```hcl
resource "hsdp_dicom_gateway_config" "dicom" {
  config_url = var.dicom_base_url
  organization_id = var.iam_org_id
  
  store_service {
    port = 104
    host_name = "foo.bar.com"
    ip_address = "1.2.3.4"
    disable_ipv6 = false
    pdu_length = 10
    artim_timeout = 20
    association_idle_timeout = 600
    network_timeout = 20
    is_secure = true
    
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

  query_service {
    port = 105
    host_name = "foo.bar.com"
    ip_address = "1.2.3.4"
    disable_ipv6 = false
    pdu_length = 10
    artim_timeout = 20
    association_idle_timeout = 600
    network_timeout = 20
    is_secure = true

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

* `config_url` - (Required) The base config URL of the DICOM Store instance
* `organization_id` - (Required) the IAM organization ID to use for authorization
* `store_service` - (Optional) Details of the CDR service account
  * `port` - Port
  * `host_name` - Host name
  * `ip_address` - IP Address
  * `disable_ipv6` - Disable IPV6
  * `pdu_length` - PDU length
  * `artim_timeout` - Artim timeout
  * `association_idle_timeout` - Association idle timeout
  * `network_timeout` - Network timeout
  * `is_secure` - Is secure
  * `application_entity` - Application entity
    * `allow_any` - Allow any
    * `ae_title` - AE title
    * `organization_id` - Organization ID
    * `service_timeout` - Service timeout
* `query_service` - (Optional) the FHIR store configuration
  * `port` - Port
  * `host_name` - Host name
  * `ip_address` - IP Address
  * `disable_ipv6` - Disable IPV6
  * `pdu_length` - PDU length
  * `artim_timeout` - Artim timeout
  * `association_idle_timeout` - Association idle timeout
  * `network_timeout` - Network timeout
  * `is_secure` - Is secure
  * `application_entity` - Application entity
    * `allow_any` - Allow any
    * `ae_title` - AE title
    * `organization_id` - Organization ID
    * `service_timeout` - Service timeout  

# Attribute reference
