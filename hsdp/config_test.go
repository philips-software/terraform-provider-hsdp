package hsdp

import (
	"testing"

	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := &Config{}

	c.Region = "us-east"
	c.Environment = "client-test"
	c.OrgAdminUsername = "foo"
	c.OrgAdminPassword = "bar"
	c.UAAPassword = "foo"
	c.UAAUsername = "bar"

	c.setupIAMClient()
	c.setupConsoleClient()

	assert.IsType(t, &iam.ErrorResponse{}, c.iamClientErr)
}
