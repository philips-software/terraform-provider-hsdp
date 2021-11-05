---
page_title: "Getting started with HSDP Terraform"
---

# Getting started with HSDP Terraform

The HSDP Terraform provider provides Lifecycle management of many HSDP resources,
including IAM, PKI, S3Creds, Edge and more. The provider is a Philips Open Source project maintained
on [Github](https://github.com/philips-software/terraform-provider-hsdp).

Support is provided through the [Github issue tracker](https://github.com/philips-software/terraform-provider-hsdp/issues)
and the `#terraform` channel on HSDP Slack.

~> The HSDP Terraform provider is not a managed service offering from HSDP, therefore please **do not open ServiceNow tickets** if you encounter issues. Instead, use one of the above-mentioned support channels. The community is pretty responsive.

## Prerequisites for using HSDP Terraform

To effectively use the HSDP Terraform provider please take into consideration the following:

- Have a HSDP LDAP account and registered SSH public key
- For services which are not self-provisioned (examples IAM, CDR): request provisioning via SNOW ticket. Your Technical Account Manager can assist you with this.
- Consider where to store your [Terraform state](https://registry.terraform.io/providers/philips-software/hsdp/latest/docs/guides/state)
