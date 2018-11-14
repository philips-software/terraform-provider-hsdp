# HSDP Terraform provider

- Website: https://www.terraform.io
- Documentation: https://github.com/philips-software/terraform-provider-hsdp/wiki
- [![Slack](https://philips-software-slackin.now.sh/badge.svg)](https://philips-software-slackin.now.sh)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

## Overview

This is a terraform provider to build/verify HSDP IAM state and other resources.
To find out more about HSDP please visit https://www.hsdp.io/

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.11 or newer (to build the provider plugin)

## Building The Provider

Clone repository somehere *outside* your $GOPATH:

```sh
$ git clone git@github.com:philips-software/terraform-provider-hsdp
$ cd terraform-provider-hsdp
$ go build .
```

Copy the resulting binary to the appropiate plugin directory e.g. `terraform.d/plugins/darwin_amd64/terraform-provider-hsdp`


## Dockerfile

A Dockerfile is provided. Example usage of the image:

```sh
$ docker build -t terraform-provider-hsdp .
$ docker run --rm -v /Location/With/Terraform/Files:/terraform -w /terraform -it terraform-provider-hsdp check
```

Automatic builds can be found on [Docker hub](https://hub.docker.com/r/philipssoftware/terraform-provider-hsdp/).

## Issues

- If you have an issue: report it on the [issue tracker](https://github.com/philips-software/terraform-provider-hsdp/issues)

## LICENSE

License is MIT
