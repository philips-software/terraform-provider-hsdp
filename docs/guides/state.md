---
page_title: "Setting up initial Terraform state"
---
# Managing Terraform state
Terraform must store state about your managed infrastructure and configuration. 
This state is used by Terraform to map real world resources to your configuration, keep track of metadata, 
and to improve performance for large infrastructures.

- [Introduction to Terraform state](https://www.terraform.io/docs/language/state/index.html)

We will discuss the various options you you have to manage Terraform state on HSDP

# S3
HSDP S3 Buckets can be used to store Terraform state. The instructions below assume some familiarity with Cloud foundry
and provisioning services using the [CF CLI](https://github.com/cloudfoundry/cli). Steps to provision an S3 Bucket:

### Log into your CF Org and space. It's advisable to create the bucket in a separate space so you can control access.

~> Make sure you restrict access to the S3 Bucket as it will contain secrets and possibly other sensitive
information that should be well protected.

### Provision a HSDP S3 Bucket
```shell
$ cf create-service hsdp-s3 s3_bucket s3-terraform
```

### Create a service key
```shell
$ cf create-service-key s3-terraform key
```

### Read out the bucket credentials
```shell
$ cf service-key s3-terraform key
```
You should see the bucket credentials on screen:
```json
{
 "api_key": "<access_key>",
 "bucket": "cf-s3-...",
 "endpoint": "s3-external-1.amazonaws.com",
 "secret_key": "<secret_key>",
 "uri": "s3://..."
}

```

### Create a `backend.tf` file
```hcl
terraform {
  backend "s3" {
    bucket = "cf-s3-..."
    key    = "project_name/some/key"
    region = "us-east-1"
  }
}
```
You can reuse a single bucket for storing multiple Terraform projects just make sure each project uses a different `key`

### Initialize the S3 backend
```shell
$ terraform init -backend-config="access_key=<api_key>" -backend-config="secret_key=<secret_key>"
```
