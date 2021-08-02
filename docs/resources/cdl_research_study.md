# hsdp_cdl_research_study

Provides a resource for managing HSDP Clinical Data Lake research studies.
A Research Study is a concept in CDL used to organize the data within Data Lake. 
It acts as a container of your clinical trial. Data will be completely isolated,
an authorization model can be enforced per Research Study.

## Example Usage

```hcl
resource "hsdp_cdl_research_study" "study_a" {
  cdl_endpoint = data.cdl_instance.cicd.endpoint
  
  title       = "Study A"
  description = "Example study A"
  study_owner = var.study_owner_id
          
  ends_at = var.ends_at
  
  data_scientist {
    user_id = data.hsdp_iam_user.scientist.id
    email   = data.hsdp_iam_user.scientist.email_address
  }
  
  uploader {
    user_id = data.hsdp_iam_user.uploaderA.id
    email   = data.hsdp_iam_user.uploaderA.email_address
  }
  
  uploader {
    user_id = data.hsdp_iam_user.uploaderB.id
    email   = data.hsdp_iam_user.uploaderB.email_address
  }
  
  monitor {
    user_id = data.hsdp_iam_user.monitor.id
    email   = data.hsdp_iam_user.monitor.email_address
  }
  
  study_manager {
    user_id = data.hsdp_iam_user.study_manager.id
    email   = data.hsdp_iam_user.study_manager.email_address
  }

  data_protected_from_deletion = false
} 
```

## Argument Reference

The following arguments are supported:

* `cdl_endpoint` - (Required) The CDL instance to query
* `title` - (Required) The name of the application
* `study_owner` - (Required, UUIDv4) The owner of the study
* `description` - (Optional) The description of the application
* `ends_at` - (Optional) The end date of the study
* `data_scientist` - (Optional) A data scientist role definition for the study
  * `user_id` (Required) The IAM user ID of the data scientist
  * `email` - (Required) The email address for this data scientist (for display purposes)
  * `institute_id` - (Optional) The institute ID associated with this role
* `monitor` - (Optional) A monitor role definition for the study
    * `user_id` (Required) The IAM user ID of the monitor
    * `email` - (Required) The email address for this monitor (for display purposes)
    * `institute_id` - (Optional) The institute ID associated with this role
* `uploader` - (Optional) An uploader role definition for the study
    * `user_id` (Required) The IAM user ID of the uploader
    * `email` - (Required) The email address for this uploader (for display purposes)
    * `institute_id` - (Optional) The institute ID associated with this role    
* `study_manager` - (Optional) An study manager role definition for the study
    * `user_id` (Required) The IAM user ID of the study manager
    * `email` - (Required) The email address for this study manager (for display purposes)
    * `institute_id` - (Optional) The institute ID associated with this role
* `data_protected_from_deletion` (Optional) Protects data from being deleted. Default is `false`


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the study

## Import

An existing research study can be imported using `terraform import hsdp_cdl_research_study`, e.g.

```shell
> terraform import hsdp_cdl_research_study.mystudy a-guid
```

