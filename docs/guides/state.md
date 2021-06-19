---
page_title: "Setting up initial Terraform state"
---
# Managing Terraform state
Terraform must store state about your managed infrastructure and configuration. 
Terraform uses the state to keep track of real world resources to your configuration, storing of metadata
and to improve performance for large infrastructure deployments.

- [Introduction to Terraform state](https://www.terraform.io/docs/language/state/index.html)

We will discuss the various options you have to manage Terraform state on HSDP

# Local
This is the default state location which Terraform uses when you do not configure a backend.
It works fine when just getting your feet wet. Terraform
will store the state in a file located in your current directory called `terraform.tfstate`. This works for small
tests and projects with only a single user. As soon as you start collaborating you will want to
use a remote state backend, like S3.

# S3
HSDP S3 Buckets can be used to store Terraform state. The instructions below assume some familiarity with Cloud foundry
and provisioning services using the [CF CLI](https://github.com/cloudfoundry/cli). Steps to provision an S3 Bucket:

### Log into your CF Org and space. 
It's advisable to create the bucket in a separate space, so you can restrict access.

!> Terraform state usually contains secrets and possibly other sensitive values related to your infrastructure and applications. Access to
state should be limited to deployment pipelines and authorized personnel only.

### Provision a HSDP S3 Bucket
It's advised to set a region
```shell
cf create-service hsdp-s3 s3_bucket s3-terraform -c '{"Region": "eu-west-1"}'
```

### Create a service key
```shell
cf create-service-key s3-terraform key
```

### Read out the bucket credentials
```shell
cf service-key s3-terraform key
```
You should see the bucket credentials on screen:
```json
{
  "api_key": "<access_key>",
  "bucket": "cf-s3-...", 
  "location_constraint": "eu-west-1",
  "endpoint": "s3-eu-west-1.amazonaws.com",
  "secret_key": "<secret_key>",
  "uri": "s3://..."
}

```

### Create a `backend.tf` file
```hcl
terraform {
  backend "s3" {
  }
}
```
You can reuse a single bucket for storing multiple Terraform projects just make sure each project uses a different `key`

### Initialize the S3 backend
Replace the values with the S3 credentails and choose a `key`
```shell
terraform init \
  -backend-config="access_key=<api_key>" \
  -backend-config="secret_key=<secret_key>" \
  -backend-config="bucket=<bucket>" \
  -backend-config="region=<region>" \
  -backend-config="key=<project_id>/<your_state_name>"
```

Your Terraform setup will now store state remotely in the S3 bucket.
