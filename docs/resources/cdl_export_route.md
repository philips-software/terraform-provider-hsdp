# hsdp_cdl_export_route

Manages HSDP Clinical Data Lake ExportRoute.

## Example Usage

```hcl
resource hsdp_cdl_export_route "expRoute1" {
  cdl_endpoint = "https://RouteFrefix-datalake.cloud.pcftest.com/store/cdl/XXXXXXX-f896-4883-80fa-5593cd69556d"
  export_route_name = "ExportTrial_for_demo55"
  description = "description11"
  display_name = "display name1"
  source_research_study {
    source_cdl_endpoint = "https://ROUTE_PREFIX-datalake.cloud.pcftest.com/store/cdl/XXXXXX-f896-4883-80fa-5593cd69556d/Study/a1467792-ef81-11eb-8ac2-477a9e3b09aa"
    allowed_dataobjects {
      resource_type = "DataObject.customDataTest123"
      associated_labels{
        label_name = "videoQualityTF7"
        approval_required = true
      } 
    }
    allowed_dataobjects {
      resource_type = "DataObject.customDataTest123"
      associated_labels{
        label_name = "videoQualityTF8"
        approval_required = true
      }
      associated_labels{
        label_name = "videoQualityTF7"
        approval_required = false
      } 
    }
  }
  auto_export = true
  destination_research_study_endpoint = "https://ROUTE_PREFIX-datalake.cloud.pcftest.com/store/cdl/XXXXXX-f896-4883-80fa-5593cd69556d/Study/5c8431e2-f4f1-11eb-bf8f-b799651c8a11"
  service_account_details {
    service_id = "SVC_ID"
    private_key = "-----BEGIN RSA PRIVATE KEY-----SVC ACC PRIVATE KEY-----END RSA PRIVATE KEY-----"
    access_token_endpoint = "https://IAM_HOST_NAME/oauth2/access_token"
    token_endpoint = "https://IAM_HOST_NAME/authorize/oauth2/token"
  }
}
```


## Argument Reference

The following arguments are supported:

* `cdl_endpoint` - (Required) The CDL instance endpoint to query
* `export_route_name` - (Required) The name of the ExportRoute
* `description` -  Description of the ExportRoute
* `display_name` -	(Required) Display Name of the ExportRoute
* `source_research_study` - (Required) Use this block to specify the details of the source Research Study
  * `source_research_study_endpoint` - (Required) "The research study endpoint of the source, for eg. https://ROUTE_PREFIX-datalake.cloud.pcftest.com/store/cdl/XXXXXX-f896-4883-80fa-5593cd69556d/Study/a1467792-ef81-11eb-8ac2-477a9e3b09aa"
  * `allowed_dataobjects` - "The DataObject details (multiple blocks of allowed_dataobjects are supported) containing resource type and the labels associated to the dataobjects" 
    * `resource_type` - (Required) "The resource type of the DataObject" 
    * `associated_labels` 
      * `label_name` - "Name of the label that is associated with the data object"
      * `approval_required` - "Boolean argument that triggers export automatically when the label is approved"
* `auto_export` - Boolean argument which shows the status of auto_export
* `destination_research_study_endpoint` - (Required) This argument represents the destination CDL endpoint 	
* `service_account_details` - (Required) This block represents the service account details
    * `service_id` - (Required) This is service_id of the service account used for the export route
    * `private_key` - (Required) The private key corresponding to the service acccount
    * `access_token_endpoint` - (Required) The access token endpoint - For ex:- "https://IAM_HOST/oauth2/access_token"
    * `token_endpoint` - (Required) The token endpoint - For ex:- "https://IAM_HOST/authorize/oauth2/token"

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the ExportRoute
* `created_by` - User who created the ExportRoute
* `created_on` - Datetime of creation of ExportRoute
* `updated_by` - The user who updated the ExportRoute
* `updated_on` - Datetime of update of ExportRoute