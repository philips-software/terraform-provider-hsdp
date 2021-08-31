# hsdp_cdl_label_definition

Manages HSDP Clinical Data Lake Label Definitions.

## Example Usage

```hcl
resource hsdp_cdl_label_definition "labeldef1" {
  cdl_endpoint    = "https://cicd-datalake.cloud.pcftest.com/store/cdl/1f5be763-f896-4883-80fa-5593cd69556d"
  study_id        = "a1467792-ef81-11eb-8ac2-477a9e3b09aa"
  label_def_name = "TF create Test4"
  description = "TF create desc"  
  label_scope = "DataObject.DICOM"
  label_name = "videoQualityTF10"
  type = "cdl/video-classification"
  labels = ["good", "bad", "acceptable", "something", "something1"]
}
```

## Argument Reference

The following arguments are supported:

* `cdl_endpoint` - (Required) The CDL instance endpoint to query
* `study_id` - (Required) The research study id under which label definition has to be created
* `label_name` - The label name
* `label_scope` - Use this parameter to specify for which CDL Data the LabelDefinition is applicable
* `labels` - Use this parameter to specify your labels, or classes. Add one label for each class.
* `label_def_name` - Name of the label definition
* `type` - Use this parameter to define the label type. Supported Values : cdl/video-classification
* `description` - Description of label definition

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the label definition
* `created_by` - User who created the label definition
* `created_on` - Timestamp the label definition was created
