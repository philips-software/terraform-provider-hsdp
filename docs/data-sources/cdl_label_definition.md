---
subcategory: "Clinical Data Lake"
---

# hsdp_cdl_label_definition

Retrieve details on HSDP Clinical Data Lake Label Definition.

## Example Usage

```hcl
data "hsdp_cdl_label_definition" "labeldef1" {
  cdl_endpoint    = "https://{{CDL-HOST}}/store/cdl/{{ORG_ID}}"
  label_def_id    = "277a5d14-86cd-4a99-92e2-7b8e898cffae"
  study_id        = "a1467792-ef81-11eb-8ac2-477a9e3b09aa"
}

output hsdp_cdl_label_definition{
  value = data.hsdp_cdl_label_definition.labeldef1
}
```

## Argument Reference

The following arguments are supported:

* `cdl_endpoint` - (Required) The CDL instance endpoint to query
* `label_def_id` - (Required) The ID of the label definition to look up
* `study_id`     - (Required) The ID of the Research Study which contains the label definition

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the Label definitions
* `label_name` -  The label name
* `label_scope` - Use this parameter to specify for which CDL Data the LabelDefinition is applicable
* `labels` - Use this parameter to specify your labels, or classes. Add one label for each class.
* `label_def_name` - Name of the label definition
* `type` - Use this parameter to define the label type. Supported Values : cdl/video-classification
* `description` - Description of label definition
* `created_by` - User who created this label definition
* `created_on` - Timestamp when label definition was created
