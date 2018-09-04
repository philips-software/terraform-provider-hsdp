# Example usage

This directory contains TF files to setup a role and a group with permissions for viewing log files (Kibana). It also contains a user definition which will create and trigger an inivitation email to IAM if the user doesn't exist.

## Files

* [hsdp.tf](hsdp.tf) - Provider configuration
* [roles.tf](roles.tf) - Role definition
* [groups.tf](groups.tf) - Group definition
* [users.tf](users.tf) - User definition
