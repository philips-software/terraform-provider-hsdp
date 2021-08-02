package hsdp

import "github.com/pkg/errors"

var (
	ErrInstanceIDMismatch       = errors.New("instanceID mismatch")
	ErrCannotCreateRootOrg      = errors.New("cannot create root organizations")
	ErrMissingParentOrgID       = errors.New("missing parent_org_id")
	ErrMissingClientID          = errors.New("missing Oauth2 client id")
	ErrMissingClientPassword    = errors.New("missing OAuth2 client password")
	ErrInvalidResponse          = errors.New("invalid response received")
	ErrResourceNotFound         = errors.New("resource not found")
	ErrIntermittent             = errors.New("intermittent error detected")
	ErrDeleteGroupFailed        = errors.New("delete group failed")
	ErrDeleteMFAPolicyFailed    = errors.New("delete of MFA policy failed")
	ErrDeleteClientFailed       = errors.New("delete client failed")
	ErrDeleteServiceFailed      = errors.New("delete service failed")
	ErrDeleteSubscriptionFailed = errors.New("delete subscription failed")
	ErrMissingOrganizationID    = errors.New("missing organization ID")
	ErrMissingIAMCredentials    = errors.New("missing IAM credentials in the hsdp provider block. Add an IAM service identity or ORG admin with proper permissions")
	ErrMissingUAACredentials    = errors.New("missing/invalid UAA credentials in the hsdp provider block")
)
