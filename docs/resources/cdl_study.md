# hsdp_cdl_study

Provides a resource for managing HSDP Clinical Data Lake Research studies

## Example Usage

```hcl
resource "hsdp_cdl_study" "study_a" {
  title = "Study A"
  description = "Example study A"
  study_owner = var.study_owner_id
          
  ends_at = var.ends_at
  
  managers =        [var.managers]
  monitors =        [var.monitors]
  uploaders =       [var.uploaders]
  data_scientists = [var.data_scientists]
} 
```

## Argument Reference

The following arguments are supported:

* `title` - (Required) The name of the application
* `study_owner` - (Required, UUIDv4) The owner of the study
* `description` - (Optional) The description of the application
* `ends_at` - (Optional) The end date of the study
* `managers` - (Optional, list(UUIDv4)) The list of study managers 
* `monitors` - (Optional, list(UUIDv4)) The list of study monitors
* `uploaders` - (Optional, list(UUIDv4)) The list of study uploaders
* `data_scientists` - (Optional, list(UUIDv4)) The list of data scientists

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the study

## Import

An existing study can be imported using `terraform import hsdp_cdl_study`, e.g.

```shell
> terraform import hsdp_cdl_study.mystudy a-guid
```

