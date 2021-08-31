# hsdp_cdl_data_type_definition

Retrieve details on HSDP Clinical Data Lake Data Type Definition.

## Example Usage

```hcl
data "hsdp_cdl_data_type_definition" "def_a" {
  cdl_endpoint = data.cdl_instance.cicd.endpoint
  dtd_id = var.dtd_id
} 

output "schema" {
  value = data.hsdp_cdl_data_type_definition.def_a.json_schema
}
```

## Argument Reference

The following arguments are supported:

* `cdl_endpoint` - (Required) The CDL instance endpoint to query
* `dtd_id` - (Required) The ID of the DTD to look up

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the DTDT
* `json_schema` -  The JSON schema describing the DTD
* `created_by` - Which entity created the DTD
* `created_on` - When the DTD was created
* `updated_by` - Which entity updated the DTD last
* `updated_on` - When the DTD was updated
