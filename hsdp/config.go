package hsdp

import (
	"github.com/hsdp/go-hsdp-iam/api"
)

type Config struct {
	api.Config
}

func (c *Config) Client() (interface{}, error) {
	client, err := api.NewClient(nil, &c.Config)
	if err != nil {
		return nil, err
	}
	err = client.Login(c.OrgAdminUsername, c.OrgAdminPassword)
	if err != nil {
		return nil, err
	}
	return client, nil
}
