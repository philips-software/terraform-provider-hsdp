# hsdp_cdl_data_type_definitions

Retrieve DTDs present in a HSDP Clinical Data Lake instance.

## Example Usage

```hcl
data "hsdp_cdl_data_type_definitions" "all" {
  cdl_endpoint = data.cdl_instance.cicd.endpoint
} 

output "all_study_titles" {
  value = data.hsdp_cdl_research_studies.all.titles
}
```

## Argument Reference

The following arguments are supported:

* `cdl_endpoint` - (Required) The CDL instance endpoint to query

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ids` -  The list of DTD GUIDs
* `names` - The names of the DTDs. This matches up with the `ids` list
