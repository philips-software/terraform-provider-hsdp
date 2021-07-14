# hsdp_cdl_instance

Retrieve details of an existing Clinical Data Lake instance (CDL.

## Example Usage

```hcl
data "hsdp_cdl_instance" "cdl" {
  base_url        = "https://my-datalake.hsdp.io"
  organization_id = var.cdl_tenant_org
}
```

## Argument Reference

The following arguments are supported:

* `base_url` - (Required) the base URL of the CDL instance. This is provided by HSDP.
* `organization_id` - (Required) the CDL tenant organization. This is provided by HSDP.

## Attributes Reference

The following attributes are exported:

* `endpoint` - The CDL store endpoint URL
