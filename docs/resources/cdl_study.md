# hsdp_cdl_study

Provides a resource for managing HSDP Clinical Data Lake Research studies.
A (Research) Study is a concept in CDL used to organize the data within Data Lake. 
It acts as a container of your clinical trial. Data will be completely isolated,
an authorization model can be enforced per Research Study.

## Example Usage

```hcl
resource "hsdp_cdl_study" "study_a" {
  title = "Study A"
  description = "Example study A"
  study_owner = var.study_owner_id
          
  ends_at = var.ends_at
  
  managers =        [var.managers]
} 
```

## Argument Reference

The following arguments are supported:

* `title` - (Required) The name of the application
* `study_owner` - (Required, UUIDv4) The owner of the study
* `description` - (Optional) The description of the application
* `ends_at` - (Optional) The end date of the study

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the study

## Import

An existing study can be imported using `terraform import hsdp_cdl_study`, e.g.

```shell
> terraform import hsdp_cdl_study.mystudy a-guid
```

