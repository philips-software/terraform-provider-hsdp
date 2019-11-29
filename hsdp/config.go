package hsdp

import (
	"github.com/philips-software/go-hsdp-api/credentials"
	"github.com/philips-software/go-hsdp-api/iam"
)

// Config contains configuration for the client
type Config struct {
	iam.Config
	S3CredsURL string

	iamClient      *iam.Client
	credsClient    *credentials.Client
	credsClientErr error
}

func (c *Config) IAMClient() *iam.Client {
	return c.iamClient
}

func (c *Config) CredsClient() (*credentials.Client, error) {
	return c.credsClient, c.credsClientErr
}

// setupIAMClient sets up an HSDP IAM client
func (c *Config) setupIAMClient() error {
	client, err := iam.NewClient(nil, &c.Config)
	if err != nil {
		return err
	}
	err = client.Login(c.OrgAdminUsername, c.OrgAdminPassword)
	if err != nil {
		return err
	}
	c.iamClient = client
	return nil
}

// setupS3CredsClient sets up an HSDP S3 Credentials client
func (c *Config) setupS3CredsClient() error {
	client, err := credentials.NewClient(c.iamClient, &credentials.Config{
		BaseURL:  c.S3CredsURL,
		Debug:    c.Debug,
		DebugLog: c.DebugLog,
	})
	if err != nil {
		c.credsClientErr = err
		return err
	}
	c.credsClient = client
	return nil
}
