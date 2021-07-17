# hsdp_cdl_research_study

Retrieve details on HSDP Clinical Data Lake research study.

## Example Usage

```hcl
data "hsdp_cdl_research_study" "study_a" {
  cdl_endpoint = data.cdl_instance.cicd.endpoint
  study_id = var.study_id
} 

output "uploaders" {
  value = data.hsdp_cdl_research_study.study_a.uploaders
}
```


## Argument Reference

The following arguments are supported:

* `cdl_endpoint` - (Required) The CDL instance endpoint to query
* `study_idr` - (Required, UUIDv4) The Research study GUID

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the study
* `description` -  The description of the application
* `study_owner` - The owner of the study
* `ends_at` - The end date of the study
* `uploaders` - The list of IAM users who have role UPLOADER
* `monitors` - The list of IAM users who have role MONITOR
* `study_managers` - The list of IAM users who have role STUDYMANAGER
* `research_scientists` - The list of IAM users who have role RESEARCHSCIENTIST

