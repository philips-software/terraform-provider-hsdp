---
page_title: "Setting up initial Terraform state"
---
# Managing Terraform state

Terraform must store state about your managed infrastructure and configuration.
Terraform uses the state to keep track of real world resources to your configuration, storing of metadata
and to improve performance for large infrastructure deployments.

- [Introduction to Terraform state](https://www.terraform.io/docs/language/state/index.html)

We will discuss the various options you have to manage Terraform state on HSDP

## Local

This is the default state location which Terraform uses when you do not configure a backend.
It works fine when just getting your feet wet. Terraform
will store the state in a file located in your current directory called `terraform.tfstate`. This works for small
tests and projects with only a single user. As soon as you start collaborating you will want to
use a remote state backend, like S3.

## S3

HSDP S3 Buckets can be used to store Terraform state. Refer to the official [Terraform S3 Backend](https://developer.hashicorp.com/terraform/language/backend/s3) documentation for setup and configuration details.

