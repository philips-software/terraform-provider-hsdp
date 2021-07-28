package hsdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := &Config{}

	c.Region = "us-east"
	c.Environment = "client-test"
	c.OrgAdminUsername = "foo"
	c.OrgAdminPassword = "bar"
	c.OAuth2ClientID = "public"
	c.UAAPassword = "foo"
	c.UAAUsername = "bar"

	c.setupIAMClient()
	c.setupConsoleClient()

	assert.NotNil(t, c.iamClientErr)
}
