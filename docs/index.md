# HSDP Provider

The HSDP provider is used to interact with HSDP REST APIs to perform adminstrative configuration of platform 
resources.

## Example Usage

```hcl
# Many variables are optional

variable "region" {}
variable "environment" {}
variable "iam_url" {}
variable "idm_url" {}
variable "oauth2_client_id" {}
variable "oauth2_password" {}
variable "org_id" {}
variable "org_admin_username" {}
variable "org_admin_password" {}
variable "shared_key" {}
variable "secret_key" {}
variable "cartel_host" {}
variable "cartel_token" {}
variable "cartel_secret" {}
variable "cartel_skip_verify" {}
variable "cartel_no_tls" {}
variable "retry_max"


## Configure the HSDP Provider

provider "hsdp" {
  region             = "us-east"
  environment        = "client-test"
  iam_url            = var.iam_url
  idm_url            = var.idm_url
  oauth2_client_id   = var.oauth2_client_id
  oauth2_password    = var.oauth2_password
  org_id             = var.org_id
  org_admin_username = var.org_admin_username
  org_admin_password = var.org_admin_password
  shared_key         = var.shared_key
  secret_key         = var.secret_key
  debug              = true
  debug_log          = "/tmp/provider.log"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The HSDP region to use [us-east, eu-west, sa1, ...]

* `environment` - (Optional) The HSDP environment to use within region [client-test, prod]

* `iam_url` - (Optional) IAM API endpoint (e.g. https://iam-client-test.us-east.philips-healthsuite.com). Auto-discovered when region and environment are specified.

* `idm_url` - (Optioanl) IDM API endpoint (e.g. https://idm-client-test.us-east.philips-healthsuite.com). Auto-discovered when region and environment are specified.

* `credentials_url` - (Optional) S3 Credenials API endpoint (e.g. https://s3creds-client-test.us-east.philips-healthsuite.com). Auto-discovered when region and environment are specified.

* `oauth2_client_id` - (Required) The OAuth2 client ID as provided by HSDP

* `oauth2_password` - (Required) The OAuth2 password as provided by HSDP

* `service_id` - (Optional) The service ID to use for IAM org admin operations (conflicts with: `org_admin_username`)

* `service_private_key` - (Optional) The service private key to use for IAM org admin operations (conflicts with: `org_admin_password`)

* `org_admin_username` - (Optional) Your IAM admin username.

* `org_admin_password` - (Optional) Your IAM admin password.

* `uaa_username` - (Optional) The HSDP CF UAA username.

* `uaa_password` - (Optional) The HSDP CF UAA password.

* `uaa_url` - (Optional) The URL of the UAA authentication service

* `org_id` - (Optional) Your IAM root ORG id as provided by HSDP

* `shared_key` - (Optional) The shared key as provided by HSDP. Actions which require API signing will not work if this value is missing.

* `secret_key` - (Optional) The secret key as provided by HSDP. Actions which require API signing will not work if this value is missing.

* `cartel_host` - (Optional) The cartel host as provided by HSDP. Auto-discovered when region and environment are specified.

* `cartel_token` - (Optional) The cartel token as provided by HSDP.

* `cartel_secret` - (Optional) The cartel secret as provided by HSDP.

* `retry_max` - (Optiona) Integer, when > 0 will use a retry-able HTTP client and retry requests when applicable.

* `debug` - (Optional) If set to true, outputs details on API calls

* `debug_log` - (Optional) If set to a path, when debug is enabled outputs details to this file

