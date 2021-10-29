---
subcategory: "Clinical Data Lake"
---

# hsdp_cdl_export_route

Retrieve details on HSDP Clinical Data Lake Export Route.

## Example Usage

```hcl
data hsdp_cdl_export_route "expRoute1" {
  cdl_endpoint = "${var.datalake_url}/store/cdl/${var.cdl_tenant_org}"
  export_route_id = "176c947c-afeb-4952-a641-80ed8878ce2d"
}
output hsdp_cdl_export_route {
  value = data.hsdp_cdl_export_route.expRoute1
}
```

## Argument Reference

The following arguments are supported:

* `cdl_endpoint` - (Required) The CDL instance endpoint to query
* `export_route_id` - (Required) The ID of the Export Route to look up

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `auto_export` - Boolean variable which indicates whether autoExport is enabled or not
* `created_by` -  The email ID of the user who created this ExportRoute
* `created_on` - Creation-Datetime of this ExportRoute
* `description` - Description string of the ExportRotue
* `destination` - Destination details of the ExportRoute (Appears as a JSON string)
* `display_name` - Display name of the ExportRoute
* `export_route_id` - The UUID of the ExportRoute
* `export_route_name` - Name given to the ExportRoute
* `source` - Source Clinical Datalake details of the ExportRoute (Appears as a JSON string)
* `updated_by` - Email ID of the user who updated the ExportRoute
* `updated_on` - Datetime when the ExportRoute was updated
