package config

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

	c.SetupIAMClient()
	c.SetupConsoleClient()

	assert.NotNil(t, c.iamClientErr)
}
