package hsdpiam

import (
	"github.com/loafoe/go-hsdpiam"
)

type Config struct {
	IAMURL               string
	IDMURL               string
	OAuth2ClientID       string
	OAuth2ClientPassword string
	OrgID                string
	OrgAdminUsername     string
	OrgAdminPassword     string
	SharedKey            string
	SecretKey            string
}

func (c *Config) Client() (interface{}, error) {
	client, err := hsdpiam.NewClient(nil, &hsdpiam.Config{
		OAuth2ClientID: c.OAuth2ClientID,
		OAuth2Secret:   c.OAuth2ClientPassword,
		SharedKey:      c.SharedKey,
		SecretKey:      c.SecretKey,
		BaseIAMURL:     c.IAMURL,
		BaseIDMURL:     c.IDMURL,
	})
	if err != nil {
		return nil, err
	}
	err = client.Login(c.OrgAdminUsername, c.OrgAdminPassword)
	if err != nil {
		return nil, err
	}
	return client, nil
}
