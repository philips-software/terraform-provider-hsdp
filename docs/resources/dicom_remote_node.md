# hsdp_dicom_remote_node
This resource manages DICOM Remote nodes

# Example Usage

```hcl
resource "hsdp_dicom_remote_node" "node1" {
  config_url = var.dicom_base_url
  organization_id = var.iam_org_two_id
  title = "Node 1"
  
  network_connection {
    port = 1000
    host_name = "foo.bar.com"
    ip_address = "1.2.3.4"
    disable_ipv6 = false
    pdu_length = 10
    artim_timeout = 20
    association_idle_timeout = 600
    network_timeout = 20
    is_secure = true
  }
}
```

# Argument reference

* `config_url` - (Required) The base config URL of the DICOM Remote node instance
* `title` - (Optional) Description of the object store
* `network_connection` - (Required) Details of the remote not network connection
  * `port` - (Required) Port
  * `host_name` - (Required) Host name
  * `ip_address` - (Optional) IP Address
  * `disable_ipv6` - (Optional) Disable IPv6
  * `pdu_length` - (Required) PDU length
  * `artim_timeout` - (Required) Artim timeout
  * `network_timeout` - (Required) Network timeout
  * `is_secure` - (Required) Secure connection
* `force_delete` - (Optional) By default remote nodes stores are not deleted by the provider (soft-delete).
  By setting this value to `true` the provider removes the remote node. We strongly suggest enabling this only for ephemeral deployments.
  
# Attribute reference
* `id` - The remote node ID

