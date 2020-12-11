# HSDP Terraform provider

- Website: https://www.hsdp.io
- Documentation: https://registry.terraform.io/providers/philips-software/hsdp/latest/docs

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

## Overview

This is a terraform provider to build/verify HSDP IAM state and other resources.
To find out more about HSDP please visit https://www.hsdp.io/

# Using the provider

**Terraform 0.13+**: To install this provider, copy and paste this code into your Terraform configuration. Then, run terraform init.

```terraform
terraform {
  required_providers {
    hsdp = {
      source = "philips-software/hsdp"
      version = ">= 0.7.4"
    }
  }
}
```

## Development requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.13.x
-	[Go](https://golang.org/doc/install) 1.15 or newer (to build the provider plugin)

### Older version

Use `v0.1.0` of this provider for Terraform 0.11 or older

## Building the provider

Clone repository somehere *outside* your $GOPATH:

```sh
$ git clone git@github.com:philips-software/terraform-provider-hsdp
$ cd terraform-provider-hsdp
$ go build .
```

Copy the resulting binary to the appropiate [plugin directory](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) e.g. `.terraform.d/plugins/darwin_amd64/terraform-provider-hsdp` 


## Dockerfile

A Dockerfile is provided. Example usage of the image:

```sh
$ docker build -t terraform-provider-hsdp .
$ docker run --rm -v /Location/With/Terraform/Files:/terraform -w /terraform -it terraform-provider-hsdp init
```

Automatic builds can be found on [Docker hub](https://hub.docker.com/r/philipssoftware/terraform-provider-hsdp/).

## Issues

- If you have an issue: report it on the [issue tracker](https://github.com/philips-software/terraform-provider-hsdp/issues)

## LICENSE

License is MIT
