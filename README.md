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
      version = ">= 0.17.0"
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
## Debugging the provider

You can build and debug the provider locally:

```sh
$ go build .
$ ./terraform-provider-hsdp -debug 
Provider started, to attach Terraform set the TF_REATTACH_PROVIDERS env var:

	TF_REATTACH_PROVIDERS='{"registry.terraform.io/philips-software/hsdp":{...}}}'
```

Copy the `TF_REATTACH_PROVIDERS` and run Terraform with this value set:

```sh
$ TF_REATTACH_PROVIDERS='...' terraform init -upgrade
$ TF_REATTACH_PROVIDERS='...' terraform plan
...
```

Terraform will now use the local running copy instead of the `philips-labs/hsdp` registry version. Happy debugging!

## Issues

* If you have an issue: report it on the [issue tracker](https://github.com/philips-software/terraform-provider-hsdp/issues)

## LICENSE

License is MIT
