package hsdp

import (
	"github.com/loafoe/go-hsdpiam"
)

type Config struct {
	hsdpiam.Config
}

func (c *Config) Client() (interface{}, error) {
	client, err := hsdpiam.NewClient(nil, &c.Config)
	if err != nil {
		return nil, err
	}
	err = client.Login(c.OrgAdminUsername, c.OrgAdminPassword)
	if err != nil {
		return nil, err
	}
	return client, nil
}
