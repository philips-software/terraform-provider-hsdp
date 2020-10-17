# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

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
