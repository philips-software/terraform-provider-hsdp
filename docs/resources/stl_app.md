# hsdp_stl_app
Manages an app resource on a STL device

## Example usage
```hcl
resource "hsdp_stl_app" "myapp" {
  name = "app.yml"
  content = file(var.myapp_yaml)
}
```

## Argument reference
* `name` - (Required) The name of the resource
* `content` - (Required) The content of the resource

## Attribute reference
* `id` - The resource ID

## Importing

An existing app resource can be imported using `terraform import hsdp_stl_app`, e.g.

```shell
> terraform import hsdp_stl_app.myapp 234
```