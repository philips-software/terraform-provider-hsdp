# HSDP Provider

Use the HSDP provider to interact with the many resources supported by [HSDP](https://www.hsdp.io). This includes amongst others many IAM entities, Container Host instances, Edge devices and even some Clinical Data Repository (CDR) resources

Use the navigation to the left to read about the available resources.

To learn the basics of Terraform using this provider, follow the hands-on [get started tutorials](https://learn.hashicorp.com/tutorials/terraform/infrastructure-as-code) on HashiCorp's Learn platform.

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
```

## Argument Reference

* `region` - (Required) The HSDP region to use [`us-east`, `eu-west`, `sa1`, `ca1`, `apac3`, ...]
* `environment` - (Optional) The HSDP environment to use within region [`client-test`, `prod`] . Default is `client-test`
* `iam_url` - (Optional) IAM API endpoint. Auto-discovered from region and environment.
* `idm_url` - (Optional) IDM API endpoint Auto-discovered from region and environment.
* `s3creds_url` - (Optional) S3 Credentials API endpoint. Auto-discovered from region and environment.
* `notification_url` - (Optional) Notification service URL
* `oauth2_client_id` - (Required) The OAuth2 client ID as provided by HSDP
* `oauth2_password` - (Required) The OAuth2 password as provided by HSDP
* `service_id` - (Optional) The service ID to use for IAM org admin operations (conflicts with: `org_admin_username`)
* `service_private_key` - (Optional) The service private key to use for IAM org admin operations (conflicts with: `org_admin_password`)
* `org_admin_username` - (Optional) Your IAM admin username.
* `org_admin_password` - (Optional) Your IAM admin password.
* `uaa_username` - (Optional) The HSDP CF UAA username.
* `uaa_password` - (Optional) The HSDP CF UAA password.
* `uaa_url` - (Optional) The URL of the UAA authentication service
* `shared_key` - (Optional) The shared key as provided by HSDP. Actions which require API signing will not work if this value is missing.
* `secret_key` - (Optional) The secret key as provided by HSDP. Actions which require API signing will not work if this value is missing.
* `cartel_host` - (Optional) The cartel host as provided by HSDP. Auto-discovered when region is specified.
* `cartel_token` - (Optional) The cartel token as provided by HSDP.
* `cartel_secret` - (Optional) The cartel secret as provided by HSDP.
* `retry_max` - (Optional) Integer, when > 0 will use a retry-able HTTP client and retry requests when applicable.
* `debug_log` - (Optional) If set to a path, when debug is enabled outputs details to this file
