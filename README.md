# HSDP Terraform provider

- Documentation on [registry.terraform.io](https://registry.terraform.io/providers/philips-software/hsdp/latest/docs)

## Overview

A Terraform provider to provision and manage state of various HSDP specific resources.
To find out more about HSDP please visit [hsdp.io](https://www.hsdp.io/)

# Using the provider

**Terraform 0.13+**: To install this provider, copy and paste this code into your Terraform configuration. Then, run terraform init.

```terraform
terraform {
  required_providers {
    hsdp = {
      source = "philips-software/hsdp"
      version = ">= 0.12.10"
    }
  }
}
```

## Development requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.14.x
-	[Go](https://golang.org/doc/install) 1.16 or newer (to build the provider plugin)

## Building the provider

Clone repository somewhere:

```sh
$ git clone git@github.com:philips-software/terraform-provider-hsdp
$ cd terraform-provider-hsdp
$ go build .
```
## Buildx Dockerfile

A Buildx Dockerfile is provided, useful for local testing. Example usage:

```sh
$ docker buildx build --push -t loafoe/terraform-provider-hsdp .
$ docker pull loafoe/terraform-provider-hsdp
$ docker run --rm -v /Location/With/Terraform/Files:/terraform -w /terraform -it loafoe/terraform-provider-hsdp init
```

## Issues

* If you have an issue: report it on the [issue tracker](https://github.com/philips-software/terraform-provider-hsdp/issues)

## LICENSE

License is MIT
