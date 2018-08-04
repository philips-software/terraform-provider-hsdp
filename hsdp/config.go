package hsdp

import (
	"github.com/philips-software/go-hsdp-api/iam"
)

// Config contains configuration for the client
type Config struct {
	iam.Config
}

// Client returns a HSDP IAM client
func (c *Config) Client() (interface{}, error) {
	client, err := iam.NewClient(nil, &c.Config)
	if err != nil {
		return nil, err
	}
	err = client.Login(c.OrgAdminUsername, c.OrgAdminPassword)
	if err != nil {
		return nil, err
	}
	return client, nil
}
