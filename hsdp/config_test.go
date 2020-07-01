package hsdp

import (
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	c := &Config{}

	c.Region = "us-east"
	c.Environment = "client-test"
	c.OrgAdminUsername = "foo"
	c.OrgAdminPassword = "bar"

	c.setupIAMClient()

	assert.Equal(t, iam.ErrNotAuthorized, c.iamClientErr)
}
