# hsdp_stl_app
Manag app resources on the device

## Example usage
```hcl
resource "hsdp_stl_app" "myapp" {
  name = "app.yml"
  content = file(var.myapp_yaml)
}
```

## Argument reference

## Attribute reference