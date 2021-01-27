# hsdp_dicom_remote_node
This resource manager DICOM remote nodes

# Example Usage
```hcl
resource "hsdp_dicom_remote_node" "node1" {
  base_url = var.dicom_base_url

  title = "Remote node 1 somewhere"
  
  port = 104
  hostname = "localhost"
  ip_address = "127.0.0.1"
  disable_ipv6 = false
  network_timeout = 3000
  is_secure = false
  
  ae_title = "The AE title of the remote"
}
```