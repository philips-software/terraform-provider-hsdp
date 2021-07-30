# hsdp_cdl_data_type_definition

Manages HSDP Clinical Data Lake Data Type Definitions.

## Example Usage

```hcl
resource "hsdp_cdl_data_type_definition" "def_a" {
  cdl_endpoint = data.cdl_instance.cicd.endpoint
  name = "my CDL schema A"
  
  json_schema = <<EOF
{
 "required": [
  "email"
 ],
 "properties": {
  "name": {
   "type": "string"
  },
  "email": {
   "type": "string"
  },
  "birthdate": {
   "type": "string"
  }
 }
}
EOF
}
```


## Argument Reference

The following arguments are supported:

* `cdl_endpoint` - (Required) The CDL instance endpoint to query
* `name` - (Required) The name of the DTD
* `description` - (Optional) The description of the DTD
* `json_schema` - (Optional) The JSON Schema describing the DTD

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the DTDT
* `created_by` - Which entity created the DTD
* `created_on` - When the DTD was created
* `updated_by` - Which entity updated the DTD last
* `updated_on` - When the DTD was updated
