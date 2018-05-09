package hsdp

import (
	"github.com/hsdp/go-hsdp-api/iam"
)

type Config struct {
	iam.Config
}

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
