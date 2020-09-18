package hsdp

import "github.com/pkg/errors"

var (
	ErrInstanceIDMismatch    = errors.New("instanceID mismatch")
	ErrNotImplementedByHSDP  = errors.New("not implemented by HSDP")
	ErrCannotCreateRootOrg   = errors.New("cannot create root orgs")
	ErrMissingParentOrgID    = errors.New("missing parent_org_id")
	ErrMissingUsername       = errors.New("missing username")
	ErrMissingPassword       = errors.New("missing password")
	ErrMissingClientID       = errors.New("missing client id")
	ErrMissingClientPassword = errors.New("missing client password")
	ErrInvalidResponse       = errors.New("invalid response received")
	ErrResourceNotFound      = errors.New("resource not found")
)
