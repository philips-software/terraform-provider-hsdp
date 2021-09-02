# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).
## v0.19.2
- Update go-hsdp-api

## v0.19.1
- DICOM: Fix for potential validation issue

## v0.19.0
- Initial AI Inference support
- NEW: Data source `hsdp_ai_inference_compute_environments`
- NEW: Data source `hsdp_ai_inference_compute_targets`
- NEW: Data source `hsdp_ai_inference_instance`
- NEW: Resource `hsdp_ai_inference_compute_environment`
- NEW: Resource `hsdp_ai_inference_compute_target`
- NEW: Resource `hsdp_ai_inference_job`
- NEW: Resource `hsdp_ai_inference_model`
- BREAKING: use `edge` instead of `stl` namespace for Edge device support
- DICOM: Fix remote node parameter reading (#109)
- DICOM: Reduce retries (#110)
- Documentation fixes

## v0.18.9
- Container Host: add readiness check

## v0.18.8
- CDL: Add export route support

## v0.18.7
- CDL: Add Label definition support
- CDL: Add 'data_protected_from_deletion' to Research Study (#97)
- PKI: Improve error handling
- IAM: Improve IAM Group deletion
- IAM: Add retry logic for email template creation
- Documentation fixes

## v0.18.6
- IAM: Add retry logic for IAM Group operations
- IAM: Better handle drift in user/service assignments in groups
- Overal improvements in error reporting (go-hsdp-api)

## v0.18.5
- CDL: Add support for Data Type Definitons

## v0.18.4
- IAM: Change variable checks. Fixes #93

## v0.18.3
- CDR: Add exponential backoff retry create with token refresh

## v0.18.2
- DICOM: Alpha quality Support for DICOM gateway configuration
- CDL: Documentation fixes
- Container Host: Fix for `keep_failed_instances` flag

## v0.18.1
- CDL: Support $grant / $revoke for data scientists, uploaders, monitors and study managers
- IAM: Workaround for IAM permissions list limitation


## v0.18.0
- Initial Clinical Data Lake (CDL) support
- NEW: Resource `hsdp_cdl_research_study`
- NEW: Data source `hsdp_cdl_instance`
- NEW: Data source `hsdp_cdl_research_study`
- NEW: Data source `hsdp_cdl_research_studies` 
- NEW: Data source `hsdp_container_host_instances`

## v0.17.2
- Upgrade go-hsdp-api
- Better honor `keep_failed_instances` for Container Host

## v0.17.1
- Update siderite and other dependencies
- Add `keep_failed_instances` attribute to Container Host resources
- Fix limit on `security_groups` on Container Host 

## v0.17.0
- The `region` is now a required argument. Environment defaults to `client-test`
- Updated documentation

## v0.16.3
- Fix refresh for PKI certs
- Fix PKI tenant update step  
- Improve error messages for hsdp_pki_cert

## v0.16.2
- Fix `alt_names` for PKI Certs

## v0.16.1
- DICOM related fixes

## v0.16.0
- Support for the HSDP Notification service

## v0.15.3
- The IAM service private_key field is now generated. This fixes some inconsistency issues

## v0.15.2
- Bring back `start_at` for `run_every` scheduling of `hsdp_function`
- Documentation fixes

## v0.15.1
- Fix ferrite backend support
- Documentation fixes

## v0.15.0
- Refactor and announce `hsdp_function` beta status
- Filter out sensitive fields from debug logs
- Add support for `ferrite` backend for `hsdp_function`
- DICOM Object stores are soft deleted by default, with option to `force_delete`

## v0.14.8
- Extra validation for `hsdp_iam_service`
- Format generated IAM Service PEM key to be more parser friendly (#72)

## v0.14.7
- [NEW] Implement `private_key` and `expires_on` configurable fields for IAM Services
- Fix `hsdp_function` start time issue

## v0.14.6
- Prevent container host cleanup for colliding hosts (#69)
- Add additional security group validation (#68)  
- Fix potential hsdp_function code collision

## v0.14.5
- Increase `volume_size` to 16000 (16T) for `hsdp_container_host` resources
- Bugfix: clean up container host instance in case of failed commands
- Documentation fixes

## v0.14.4
- [NEW] `cron` support for `hsdp_function.schedule` configuration
- [NEW] `timeout` support `for hsdp_function.schedule` configuration
- Fix duplicate debug logging output

## v0.14.3
- Support CDR Org delete with optional support for $purge
- Add support for `image` field for `hsdp_container_host`
- Description fields for IAM groups and roles are now optional

## v0.14.2
- Fix some DICOM optional fields
- Improve endpoint auto-discovery
- Work on guides

## v0.14.1
- Update S3Policy actions list

## v0.14.0
- [NEW] hsdp_function resource

## v0.13.5
- Fix state issue in DICOM repository

## v0.13.4
- Improve DICOM repository onboarding

## v0.13.3
- [NEW] Support for setting permissions, owner and group for CH files
- Detect copy errors for SSH copy files

## v0.13.2
- Fix clear_on_destroy state

## v0.13.1
- Documentation fix

## v0.13.0
- Support `ensure_tcp` and `ensure_udp` in STL firewall exception config
- Documentation fixes

## v0.12.12
- Workaround for IAM profile update issue

## v0.12.11
- Documentation fixes
- Fix hsdp_iam_user.mobile field updating
- Fix corner case where IAM generates error 104 on profile update

## v0.12.10
- [NEW] optional `password` argument for immediate activation of `hsdp_iam_user`
- Support `hsdp_iam_user` field updates (first_name, last_name, login, email)

## v0.12.9
- Suppress global_reference_id diffs changes when generated

## v0.12.8
- Better error reporting and fix root cause of crashing bug

## v0.12.7
- Make global_reference_id optional for Application and Proposition

## v0.12.6
- Fix crashing bug in create IAM application

## v0.12.5
- Add retry code to overcome IAM race condition in certain situations

## v0.12.4
- Fix DICOM onboarding when provisioning IAM groups during the same run

## v0.12.3
- Add missing fields for DICOM
- Minor documentation fixes

## v0.12.2
- Fix STL cert update issue 
- Remove last_update fields as it produced inconsistent state

## v0.12.1
- Improve Proposition and Application resource lifecycle and error handling 
- Sync STL resources by default now. Users can choose to batch this using `hsdp_stl_sync`

## v0.12.0
- [NEW] Secure Transport Layer (STL) support to manage Edge devices

## v0.11.3
- Fix default IAM OAuth2 client TTLs

## v0.11.2
- Fix documentation

## v0.11.0
- NEW: HSDP PKI initial support

## v0.10.0
- NEW: DICOM config support

## v0.9.4
- Better cleanup logic for failed container host provisions

## v0.9.3
- Recovery code for Cartel HTTP 500 error during create

## v0.9.2
- Improve error handling for Cartel

## v0.9.1
- Add support for file sources
- Improve error handling

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
