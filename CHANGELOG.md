# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## v0.38.11

- MDM: expose MDM proposition guid #272

## v0.38.10

- Core: support credentials config via JSON #270

## v0.38.9

- PKI: make UAA credentials fully optional #267

## v0.38.8

- IAM: conditionally check IAM Device/User mixups #265

## v0.38.7

- PKI: support for region sa1
- PKI: add triggers field

## v0.38.6

- Chore: use generated password everywhere #258

## v0.38.5

- IAM Introspect: add scopes field
- CDR: Improve the CDR Org example
- Core: Update Google FHIR #255
- Core: fix retry logic #257

## v0.38.4

- Autoscaler: handle state drift

## v0.38.3

- IAM Role: wrap delete in retry block
- IAM Email template: wrap delete call in retry
- Notification: fix documentation 

## v0.38.2

- Core: better handling of 502 BAD Gateway responses

## v0.38.1

- IAM Service: fix migration

## v0.38.0

- IAM Device: add support for IAM devices (#244)
- IAM Group: allow device assignments
- Container Host: improve error handling (#236)
- Documentation: update Container Host guide

## v0.37.6

- PKI: update config for real

## v0.37.5

- PKI: add eu-west production
- IAM: handle intermittent issues with Group updates

## v0.37.4

- Discovery: add hsdp_discovery_service resource

## v0.37.3

- DICOM Repository: check for pre existing repo and report
- DICOM Notification: deprecate/remove organization_id field

## v0.37.2

- Core: Allow "production" as environment
- DICOM: Add store_as_composite argument (#230)
- IAM: Bump go-hsdp-api version to fix IAM consistency issue

## v0.37.1

- Notification: fix non principal authentication (#227)

## v0.37.0

- Deprecation WARNING: hsdp_container_host_exec is being deprecated. It will be removed in v0.40.0
- Edge: add support for Principal blocks
- MDM: workaround for refreshing issues with standard service references (#225)

## v0.36.0

- IAM: Add support for Role Sharing Policy (#220)

## v0.35.4

- IAM Service: add self_managed_expires_on field for better Crossplane compatibility (potentially breaking change!)

## v0.35.3

- IAM Client: fix drift detection (#218)

## v0.35.2

- IAM Client: allow description change without recreate. Add acceptance test

## v0.35.1

- IAM: relax service requirement field

## v0.35.0

- CDR: New Practitioner resource
- Notification: support Principal block to use different IAM Service accounts

## v0.34.4

- MDM: fix crashing bug

## v0.34.3

- CDR: typo fix in doc

## v0.34.2

- CDR: Fix CDR subscription documentation and status field

## v0.34.1

- MDM: Proposition state drift detection fix

## v0.34.0

- CDR: refactor hsdp_cdr_fhir_store data source, including '/store/fhir' is mandatory but not breaking
- CDR: add basic acceptance tests
- IAM: Proposition state drift detection fix

## v0.33.6

- CDR: fix subscription read for STU3

## v0.33.5

- IAM: detect missing OAuth2 credentials
- IAM: update hsdp_iam_role docs (#204)
- Docker: report warning in case GraphQL query for all repositories fails
- CDR: fix subscription example
- IAM: add additional fields to introspect

## v0.33.4

- MDM: remove extra read in Proposition create setup
- MDM: add more acceptances tests

## v0.33.3

- Docker: add latest_tag and full_name attributes (#201)

## v0.33.2

- IAM: fix introspect data source bug
- IAM: add more acceptance tests

## v0.33.1

- IAM: update hsdp_iam_permissions data source
- IAM: add hsdp_iam_permission data source

## v0.33.0

- IAM: workaround for HTTP 403/403 issue (#199)
- IAM: add preliminary acceptances tests

## v0.32.2

- Doc: fix a typo

## v0.32.1

- IAM: support for synchronous IAM organization deletes
- CDR: fix subscription importing (#188)

## v0.32.0

- IAM: various fixes for handling with managing_organization changes
- MDM: support for Bootstrap client scope management

## v0.31.3

- IAM Role: better handling of invalid permissions (#190)

## v0.31.2

- CDR Org: fix importing of CDR ORGs

## v0.31.1

- Connect MDM: Fix standard service data source (#185)

## v0.31.0

- IAM Service: Support new token_validity argument (#184)

## v0.30.13

- IAM: Fix for CDKTF support

## v0.30.12

- Container Host: bump hsdp_container_host_instances version (#180)

## v0.30.11

- Container Host: new Container Host data source (#180)
- Container Host: fix fields in hsdp_container_host_instances data source (#180)

## v0.30.10

- CDR: better handling of soft-deleted ORGs (#175)

## v0.30.9

- IAM: prevent uninitialized IAM client from being used (#173)

## v0.30.8

- CDR ORG: fix crashing bug (#173)

## v0.30.7

- CDR Subscription: fix crashing bug (#173)

## v0.30.6

- Connect MDM: documentation fix

## v0.30.5

- Connect MDM: optionally bypass MDM scope limitations

## v0.30.4

- Connect MDM: Add ServiceAgent data resources

## v0.30.3

- IAM: Fix user fetch calls in IAM Grousp (#170)
- Connect MDM: Fix empty secrets issue (#169)

## v0.30.2

- Notification: support querying for topic by name (#165)

## v0.30.1

- Connect MDM: workaround for StandardService query issue

## v0.30.0

- IAM: improve token handling

## v0.29.9

- IAM: Add hsdp_iam_token data source (#167)

## v0.29.8

- Notification: fix data source IDs

## v0.29.7

- Container Host: Fix add/remove user groups API calls
- Container Host: Initial guide published
- Container Host: Add security group data sources

## v0.29.6

- Connect MDM: add IAM application GUID attribtes

## v0.29.5

- DICOM Config: fix description field handling (#164)

## v0.29.4

- DICOM Notification: fix destroy call (#160)

## v0.29.3

- Connect MDM: add DataType datasource

## v0.29.2

- Connect MDM: Add MDM Bucket data source

## v0.29.1

- Connect MDM: expose additional data fields in hsdp_connect_mdm_application

## v0.29.0

- Core: add support for apac2 (Tokyo) region

## v0.28.5

- DICOM Notification: add workaround for broken API design
- Container Host: try to recover from flaky salt failures

## v0.28.4

- Core: upgrade go-hsdp-api

## v0.28.3

- DICOM: Fix JSON issue

## v0.28.2

- DICOM: Fix notification implementation

## v0.28.1

- DICOM: Fix crashing bug

## v0.28.0

- DICOM: Add new notification feature (December 2021 reelase)

## v0.27.14

- Container Host: fix keep_failed_instances notifications
- Core: upgrade Terraform plugin SDK to v2.10.0

## v0.27.13

- IAM: Fix documentation link

## v0.27.12

- IAM: Add hsdp_iam_email_templates data source (#155) 

## v0.27.11

- IAM: Add hsdp_iam_group_membership (#152)

## v0.27.10

- Edge: fix crashing bug and app state handling

## v0.27.9

- MDM: add retry calls to read operations as well. Fixes on-the-fly permission assignment runs

## v0.27.8

- MDM: minor performance improvement in resource creation

## v0.27.7

- MDM: show more details in HTTP 422 flow handling

## v0.27.6

- MDM Application: support description and default_group_guid arguments
- Docker Service Key: support drift detection
- IAM Group: display server error message
- Documentation: update MDM docs
- Core: update go-hsdp-api
- Function: more error message improvements
- Core: more informative error message on missing credentials

## v0.27.5

- IAM User: fix heuristics for auto importing existing users

## v0.27.4

- IAM User: fix create user handling if the account already exists
- DICOM Gateway: retry read calls

## v0.27.3

- IAM Group: refactor group handling (#147)
- Docs: fix broken MDM bucket example
- Core: Upgrade terraform-plugin-sdk to v.2.9.0

## v0.27.2

- Function: update instructions for siderite v0.12.0

## v0.27.1

- MDM: Documentation fixes
- MDM: Add hsdp_connect_mdm_firmware_distribution_request

## v0.27.0

- NEW: Connect MDM support
- Documentation: example fixes

## v0.26.6

- IAM Client: fix consent implied handling
- IAM Client: add data source

## v0.26.5

- Core: upgrade go-hsdp-api
- Dcoumentation: updates

## v0.26.4

- Core: upgrade go-hsdp-api
- Documentation: typo fixes

## v0.26.3

- Documentation: move hspd_function guide to right subcategory

## v0.26.2

- Documentation: add subcategories

## v0.26.1

- IAM Email Templates: handle server side case changes

## v0.26.0

- NEW HSDP Docker Registry support: manage namespaces and repositories

## v0.25.2

- IAM: Add retry logic to additional IAM resources

## v0.25.1

- Guide: fix S3 backend command
- Provider: update go-hsdp-api

## v0.25.0

- CDR: Support STU3 and R4 FHIR resources. Defaults to STU3, no change to existing resources
- Provider: Add validation to region and environment fields (#54)
- IAM Service: Mark expires_on as computed (#94)

## v0.24.1

- Made `region` optional and default to `us-east`
- AI Workspace: fix reading bug
- Test: skeleton code

## v0.24.0

- Chore: massive refactoring of package namespace

## v0.23.3

- PKI: Fix schema bug

## v0.23.2

- Fix hsdp_iam_user data source

## v0.23.1

- NEW: IAM Users data source: `data.hsdp_iam_users` 
- NEW: IAM Email Activation resource: `hsdp_iam_email_activation`
- Container Host: improve commands error reporting

## v0.23.0

- NEW: IAM SMS Gateway configuration support: `hsdp_iam_sms_gateway`
- NEW: IAM SMS Templates configuration: `hsdp_iam_sms_template`
- NEW: Support provider credentials and settings from the Environment
- NEW: IAM User resources supports setting preferred language and communication channel
- CDL: Fix study conflict resolution

## v0.22.2

- DICOM: Fix unexpected recreate of dicom_object_store due to API changes

## v0.22.1

- DICOM: Add query param (#125)
- DICOM: Fix hsdp_dicom_store_config hash resources
- PKI: Fix hash resources
- Edge: Fix hash resources
- Autoscaler: Fix hash resources
- DICOM: Fix hash resources

## v0.22.0

- DICOM Gateway: Breaking change: new 'organization_id' required field
- DICOM: Add proper Hash functions for nested resources

## v0.21.6

- Container Host: user is optional
- Container Host: add additional checks and fix order

## v0.21.5

- IAM: [service] remove self-managed certificate, it's an anti-pattern

## v0.21.4

- IAM: [service] clear private key when self-managed credentials are used
- IAM: read after create improvements
- Container Host: bump number of retries container host ready check
- Container Host: credentials validation check before provisioning
- Config: fix 'sliding_expired_on' value

## v0.21.3

- DICOM: Ensure ForceNew is pervasive for remote nodes

## v0.21.2

- DICOM: Use different type structures for certain API endpoints

## v0.21.1

- IAM: Ignore case for login and email fields
- Container Host: documentation fixes


## v0.21.0

- IAM: hsdp_iam_group and hsdp_iam_role data sources (#122)
- Function: propagate timeout to Iron tasks
- Container Host: support capturing output from commands (#120)

## v0.20.8

- DICOM: Fix JSON field names

## v0.20.7

- DICOM: Fix JSON rendering issue

## v0.20.6

- IAM: Fix issue with self_managed_key
- IAM: Fix perma-diff when changing Org names
- CDR: Handle Subscription drift detection
- Function: update siderite-backend version
- DICOM: Fix crashing bug

## v0.20.5

- DICOM Gateway: fix refresh and destroy for config resource

## v0.20.4

- DICOM Gateway: various fixes based on API changes
- Container Host: support for SSH-agent authentication

## v0.20.3

- DICOM Gateway: don't propagate secure toggle field

## v0.20.2

- DICOM Gateway: remove unused field

## v0.20.1

- AI: More consistent naming convention for service instances
- DICOM Gateway: use pointers in structs to satisfy validations

## v0.20.0

- Initial AI Workspace support
- NEW: Data source `hsdp_ai_workspace_compute_targets`
- NEW: Data source `hsdp_ai_workspace`
- NEW: Resource `hsdp_ai_workspace_compute_target`
- NEW: Resource `hsdp_ai_workspace`
- DICOM Gateway: fix more field reads

## v0.19.10

- DICOM Gateway: Fix various structures

## v0.19.9

- IAM: Fix detection of purged user accounts

## v0.19.8

- IAM: Do not error out in case IAM user is not found using data source

## v0.19.7

- IAM: Proper error reporting in case of missing CLIENT.SCOPE permissions
- DICOM Gateway: add title and description fields

## v0.19.6

- Config: improve documentation (#106)
- Container Host: increase command limit to 50

## v0.19.5

- IAM: Fix `application_id` changes on IAM Service identities

## v0.19.4

- Expose `service_id` and `org_admin_username` through `hsdp_config` (#113)

## v0.19.3

- Fix authentication issue when using service identities

## v0.19.2

- Update go-hsdp-api

## v0.19.1

- DICOM: Fix for potential validation issue

## v0.19.0

- Initial AI Inference support
- NEW: Data source `hsdp_ai_inference_compute_environments`
- NEW: Data source `hsdp_ai_inference_compute_targets`
- NEW: Data source `hsdp_ai_inference_service_instance`
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

- CDL: Add support for Data Type Definitions

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
