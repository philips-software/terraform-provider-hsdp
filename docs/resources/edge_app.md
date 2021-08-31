# hsdp_edge_app

Manages an app resource on an Edge device. At this time resources are synced immediately to the device after create or update.

## Example usage

```hcl
resource "hsdp_edge_app" "myapp" {
  serial_number = var.serial_number
  
  name = "app.yml"
  content = file(var.myapp_yaml)
}
```

## Argument reference

* `serial_numbe` - (Required) The serial number of the device to deploy this app resource on
* `name` - (Required) The name of the resource
* `content` - (Required) The content of the resource
* `sync` - (Optional, boolean) Sync the resource after mutation. Current default behaviour at system level is to sync immediately, but this might change in future updates.

## Attribute reference

* `id` - The resource ID

## Importing

An existing app resource can be imported using `terraform import hsdp_edge_app`, e.g.

```shell
terraform import hsdp_edge_app.myapp 234
```
