# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## v0.9.0
- [NEW] Use SSH file and commands directies for Container Host
- [NEW] Add hsdp_container_host_exec as replacement for null_resource

## v0.8.9
- Add proxy support for Cartel connections

## v0.8.8
- [NEW] Added hsdp_iam_email_template to manage IAM custom email templates

## v0.8.7
- Validate Container Host tags

## v0.8.6
- Support for setting subnet for Container Hosts
- Fix Container Host import support
- Update Terraform to 0.14.4

## v0.8.5
- Wrap more error conditions
- Use UTC timezone for FHIR parsing

## v0.8.4
- Add additional error messages

## v0.8.3
- Fix documentation

## v0.8.2
- Refactor CDR resource naming after some trial use
- Add `part_of` attribute to `hsdp_cdr_org`

## v0.8.1
- NEW: Add hsdp_cdr_subscription

## v0.8.0
- NEW: Clinical Data Repository (CDR) onboarding support

## v0.7.8
- Handle missing Role delete capability of IAM gracefully
- Fix crashing bug

## v0.7.7
- go-hsdp-api bugfix in the console API client

## V0.7.6
- Improve autoscaler support
- Fix documentation

## v0.7.5
- NEW: Add data source hsdp_container_host_subnet_types
- Container Host: add subnet_type configuration (public, private)

## v0.7.4
- Add validation checks and update documentatin for Container Host

## v0.7.3
- Implement data.hsdp_iam_service

## v0.7.2
- Use legacy fallback for data.hsdp_user

## v0.7.1
- Use Go 1.15.5
- Fix linting errors

## v0.7.0
- Upgrade to Terraform v2 SDKs

## v0.6.8
- Update to latest v1 SDKs

## v0.6.7
- Increase default timeouts for Container Host
- Fix documentationt

## v0.6.6
- NEW: Support for tagging Container Host instances
- Improve error handling for Container Host provisioning
- Fix Dockerfile for local development

## v0.6.5
- NEW: hsdp_iam_password_policy

## v0.6.4
- Add domain to hsdp config data source
- Documentation fixes

## v0.6.3
- NEW: hsdp_iam_application data source
- Fix data ID for hsdp_iam_proposition

## v0.6.2
- NEW: hsdp_iam_proposition data source
- Updated hsdp_iam_org resource to include additional fields
- Implement hsdp_iam_org deletion

## v0.6.1
- Bugfix release

## v0.6.0
- NEW: hsdp_metrics_autoscaler resource
- Migrate Terraform PLugin SDK
- Upgrade to Terraform 0.13.1

## v0.5.0
- NEW: hsdp_container_host
- Migrate to Terraform Plugin SDK
- Handle externally deleted resources
- Upgrade to Terraform 0.12.25

## v0.4.0
- Switch user API to v2 (breaking change!)
- New user login field
- Support user deletion
- Shared key and secret and now optional
- Upgrade to Terraform 0.12.24
- Support adding service identities to groups

## v0.3.0
- Upgrade to Terraform 0.12.23

## v0.2.0

- Upgrade to Terraform 0.12.x

## v0.1.0

- Initial implementation
- Application (CRUD)
- Client (CRUD)
- Group (CRUD)
- Organization (CRU)
- Permission (CRUD)
- Proposition (CRUD)
- Role (CRUD)
- User (CR)

[Unreleased]: https://github.com/philips-software/terraform-provider-hsdp/compare/9b82310...HEAD
