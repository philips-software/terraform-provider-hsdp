# HSDP Terraform provider

## Overview

This is a terraform provider to build/verify HSDP IAM state and other resources

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.10 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/philips-software/terraform-provider-hsdp`

```sh
$ mkdir -p $GOPATH/src/github.com/philips-software; cd $GOPATH/src/github.com/philips-software
$ git clone git@github.com:philips-software/terraform-provider-hsdp
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/philips-software/terraform-provider-hsdp
$ go build .
```

Copy the binary to the appropiate plugin directory e.g. `terraform.d/plugins/darwin_amd64/terraform-provider-hsdp`

## Issues

- If you have an issue: report it on the [issue tracker](https://github.com/philips-software/terraform-provider-hsdp/issues)

## LICENSE

License is MIT
