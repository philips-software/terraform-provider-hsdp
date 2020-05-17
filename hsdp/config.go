package hsdp

import (
	"github.com/philips-software/go-hsdp-api/cartel"
	"github.com/philips-software/go-hsdp-api/credentials"
	"github.com/philips-software/go-hsdp-api/iam"
)

// Config contains configuration for the client
type Config struct {
	iam.Config
	S3CredsURL       string
	CartelHost       string
	CartelToken      string
	CartelSecret     string
	CartelNoTLS      bool
	CartelSkipVerify bool

	iamClient       *iam.Client
	cartelClient    *cartel.Client
	credsClient     *credentials.Client
	credsClientErr  error
	cartelClientErr error
}

func (c *Config) IAMClient() *iam.Client {
	return c.iamClient
}

func (c *Config) CartelClient() (*cartel.Client, error) {
	return c.cartelClient, c.cartelClientErr
}

func (c *Config) CredentialsClient() (*credentials.Client, error) {
	return c.credsClient, c.credsClientErr
}

func (c *Config) CredentialsClientWithLogin(username, password string) (*credentials.Client, error) {
	newIAMClient, err := c.iamClient.WithLogin(username, password)
	if err != nil {
		return nil, err
	}
	return credentials.NewClient(newIAMClient, &credentials.Config{
		BaseURL:  c.S3CredsURL,
		Debug:    c.Debug,
		DebugLog: c.DebugLog,
	})
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
func (c *Config) setupS3CredsClient() {
	client, err := credentials.NewClient(c.iamClient, &credentials.Config{
		BaseURL:  c.S3CredsURL,
		Debug:    c.Debug,
		DebugLog: c.DebugLog,
	})
	if err != nil {
		c.credsClient = nil
		c.credsClientErr = err
		return
	}
	c.credsClient = client
	return
}

// setupCartelClient sets up an Cartel client
func (c *Config) setupCartelClient() {
	client, err := cartel.NewClient(nil, cartel.Config{
		Host:       c.CartelHost,
		Token:      c.CartelToken,
		Secret:     []byte(c.CartelSecret),
		NoTLS:      c.CartelNoTLS,
		SkipVerify: c.CartelSkipVerify,
		Debug:      c.Debug,
		DebugLog:   c.DebugLog,
	})

	if err != nil {
		c.cartelClient = nil
		c.cartelClientErr = err
		return
	}
	c.cartelClient = client
	return
}
