# HSDP Provider

Use the HSDP provider to interact with the many resources supported by [HSDP](https://www.hsdp.io). This includes amongst others IAM entities, Container Host instances, Edge devices and even some Clinical Data Repository (CDR) resources

Use the navigation to the left to read about the available resources.

To learn the basics of Terraform, follow the hands-on [get started tutorials](https://learn.hashicorp.com/tutorials/terraform/infrastructure-as-code) on HashiCorp's Learn platform.

## Example usage

```hcl
provider "hsdp" {
  region             = "us-east"
  environment        = "client-test"
  oauth2_client_id   = var.oauth2_client_id
  oauth2_password    = var.oauth2_password
  org_admin_username = var.org_admin_username
  org_admin_password = var.org_admin_password
}

resource "hsdp_iam_org" "hospital_a" {
  name          = "HOSPITAL_A"
  description   = "HOSPITAL A"
  parent_org_id = data.hsdp_iam_org.root.id
}
```

## Authentication

The HSDP provider can read credentials and settings from the Environment or as
arguments in its provider block. The following environment variables are recognized

| Environment                  | Argument            | Required | Default     |
|------------------------------|---------------------|----------|-------------|
| HSDP_REGION                  | region              | Optional | us-east     |
| HSDP_ENVIRONMENT             | environment         | Optional | client-test |
| HSDP_CARTEL_HOST             | cartel_host         | Optional |             |
| HSDP_CARTEL_SECRET           | cartel_secret       | Optional |             |
| HSDP_CARTEL_TOKEN            | cartel_token        | Optional |             |
| HSDP_IAM_SERVICE_ID          | service_id          | Optional |             |
| HSDP_IAM_SERVICE_PRIVATE_KEY | service_private_key | Optional |             |
| HSDP_IAM_ORG_ADMIN_USERNAME  | org_admin_username  | Optional |             |
| HSDP_IAM_ORG_ADMIN_PASSWORD  | org_admin_password  | Optional |             |
| HSDP_IAM_OAUTH2_CLIENT_ID    | oauth2_client_id    | Optional |             |
| HSDP_IAM_OAUTH2_PASSWORD     | oauth2_password     | Optional |             |
| HSDP_SHARED_KEY              | shared_key          | Optional |             |
| HSDP_SECRET_KEY              | secret_key          | Optional |             |
| HSDP_UAA_USERNAME            | uaa_username        | Optional |             |
| HSDP_UAA_PASSWORD            | uaa_password        | Optional |             |
| HSDP_DEBUG_LOG               | debug_log           | Optional |             |
| HSDP_DEBUG_STDERR            | debug_stderr        | Optional |             |

## Argument Reference

In addition to generic provider arguments (e.g. alias and version), the following arguments are supported in the HSDP provider block:

* `region` - (Required) The HSDP region to use [`us-east`, `eu-west`, `sa1`, `ca1`, `apac3`, ...]. Default is `us-east`
* `environment` - (Optional) The HSDP environment to use within region [`client-test`, `prod`] . Default is `client-test`
* `credentials` - (Optional) Can point to a JSON file containing values for all fields here
* `iam_url` - (Optional) IAM API endpoint. Auto-discovered from region and environment.
* `idm_url` - (Optional) IDM API endpoint Auto-discovered from region and environment.
* `s3creds_url` - (Optional) S3 Credentials API endpoint. Auto-discovered from region and environment.
* `notification_url` - (Optional) Notification service URL. Auto-discovered from region and environment.
* `oauth2_client_id` - (Optional) The OAuth2 client ID as provided by HSDP
* `oauth2_password` - (Optional) The OAuth2 password as provided by HSDP
* `service_id` - (Optional) The service ID to use for IAM org admin operations (conflicts with: `org_admin_username`)
* `service_private_key` - (Optional) The service private key to use for IAM org admin operations (conflicts with: `org_admin_password`)
* `org_admin_username` - (Optional) Your IAM admin username.
* `org_admin_password` - (Optional) Your IAM admin password.
* `uaa_username` - (Optional) The HSDP CF UAA username.
* `uaa_password` - (Optional) The HSDP CF UAA password.
* `uaa_url` - (Optional) The URL of the UAA authentication service. Auto-discovered from region.
* `mdm_url` - (Optional) The base URL of the MDM service. Auto-discovered from region and environment.
* `shared_key` - (Optional) The shared key as provided by HSDP. Actions which require API signing will not work if this value is missing.
* `secret_key` - (Optional) The secret key as provided by HSDP. Actions which require API signing will not work if this value is missing.
* `cartel_host` - (Optional) The cartel host as provided by HSDP. Auto-discovered from region.
* `cartel_token` - (Optional) The cartel token as provided by HSDP.
* `cartel_secret` - (Optional) The cartel secret as provided by HSDP.
* `retry_max` - (Optional) Integer, when > 0 will use a retry-able HTTP client and retry requests when applicable.
* `debug_log` - (Optional) If set to a path, when debug is enabled outputs details to this file
* `debug_stderr` - (Optional) If set to true sends debug logs to `stderr`
