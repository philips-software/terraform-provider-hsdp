# hsdp_dicom_remote_node
This resource manages DICOM Remote nodes using HSDP provisioned DICOM Store configuration service.

# Example Usage

```hcl
resource "hsdp_dicom_remote_node" "remote_node_1" {
  config_url = var.dicom_base_url
  organization_id = var.site_org
 
  title = "Node 1"
  ae_title = "AeTitelNode1" 
  
  network_connection {
    port = 1000
    hostname = "foo.bar.com"
    ip_address = "1.2.3.4"
    disable_ipv6 = false
    pdu_length = 10
    artim_timeout = 20
    association_idle_timeout = 600
    network_timeout = 20
    is_secure = false
  }
}
```

# Argument reference

* `config_url` - (Required) The base config URL of the DICOM Store
* `site_organization_id` - (Required) Site Organization ID for which Gateway to be deployed
* `title` - (Optional) Remote Noe Description
* `ae_title` - (Required) Remote Node AETtile. Allowed characters for aetitle are `A-Za-z0-9\\s/+=_-`
* `network_connection` - (Required) Details of the remote not network connection
  * `port` - (Required) Remte Node Port
  * `hostname` - (Required) Remote Node Host name
  * `ip_address` - (Optional) Remote Nde IP Address
  * `disable_ipv6` - (Optional) Disable IPv6. Default `false`
  * `pdu_length` - (optional) PDU length. Default `65535`
  * `artim_timeout` - (optional) Time-out waiting for A-ASSOCIATE RQ PDU on open TCP/IP connection. Artim timeout. Default `3000 ms`
  * `associationIdleTimeOut` - (Optional) Association Idle Timeout. Default `4500 ms`
  * `network_timeout` - (optional) Network timeout. Default `3000 ms`
  * `is_secure` - (Required) Secure connection. Boolean `true` or `false`. Default `false`
* `force_delete` - (Optional) By default remote nodes are not deleted by the provider (soft-delete).
  By setting this value to `true` the provider removes the remote node. We strongly suggest enabling this only for ephemeral deployments.
  
# Attribute reference
* `id` - The remote node ID

