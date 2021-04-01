---
page_title: "Setting up initial Terraform state"
---
# Managing Terraform state
Terraform must store state about your managed infrastructure and configuration. 
This state is used by Terraform to map real world resources to your configuration, keep track of metadata, 
and to improve performance for large infrastructures.

Please go through the articles below for further details on state:
- [Introduction to State](https://www.terraform.io/docs/language/state/index.html)
- [Purpase of Terraform State](https://www.terraform.io/docs/language/state/purpose.html)

The following sections talk about the various backends you can use to store your state

# S3
HSDP S3 Buckets can be used to store Terraform state. The instructions below assume some familiarity with Cloud foundry
and provisioning services using the [CF CLI](https://github.com/cloudfoundry/cli). Steps to provision an S3 Bucket:

1. Log into your CF Org and space. It's advisable to create the bucket in a separate space so you can control access.

~> Important: make sure you restrict access to the S3 Bucket as it will contain secrets and possibly other sensitive
information that should be well protected.

2. Provision a HSDP S3 Bucket
```shell
$ cf create-service hsdp-s3 s3_bucket s3-terraform
```

3. Create a service key
```shell
$ cf create-service-key s3-terraform key
```

4. Read out the bucket credentials
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

5. Create a `backend.tf` file:
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

6. Initialize the S3 backend:
```shell
$ terraform init -backend-config="access_key=<api_key>" -backend-config="secret_key=<secret_key>"
```
