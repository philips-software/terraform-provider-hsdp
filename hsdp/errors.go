package hsdp

import "github.com/pkg/errors"

var (
	ErrInstanceIDMismatch   = errors.New("instanceID mismatch")
	ErrNotImplementedByHSDP = errors.New("not implemented by HSDP")
	ErrCannotCreateRootOrg  = errors.New("cannot create root orgs")
	ErrMissingParentOrgID   = errors.New("missing parent_org_id")
)
