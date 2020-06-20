package hsdp

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/philips-software/go-hsdp-api/cartel"
	"github.com/philips-software/go-hsdp-api/credentials"
	"github.com/philips-software/go-hsdp-api/iam"
	"net/http"
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
	RetryMax         int

	iamClient       *iam.Client
	cartelClient    *cartel.Client
	credsClient     *credentials.Client
	credsClientErr  error
	cartelClientErr error
	iamClientErr    error
}

func (c *Config) IAMClient() (*iam.Client, error) {
	return c.iamClient, c.iamClientErr
}

func (c *Config) CartelClient() (*cartel.Client, error) {
	return c.cartelClient, c.cartelClientErr
}

func (c *Config) CredentialsClient() (*credentials.Client, error) {
	return c.credsClient, c.credsClientErr
}

func (c *Config) CredentialsClientWithLogin(username, password string) (*credentials.Client, error) {
	if c.iamClientErr != nil {
		return nil, c.iamClientErr
	}
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
func (c *Config) setupIAMClient() {
	standardClient := http.DefaultClient
	if c.RetryMax > 0 {
		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = 5
		standardClient = retryClient.StandardClient()
	}
	c.iamClient = nil
	client, err := iam.NewClient(standardClient, &c.Config)
	if err != nil {
		c.iamClientErr = err
		return
	}
	if c.OrgAdminUsername == "" {
		c.iamClientErr = ErrMissingUsername
		return
	}
	if c.OrgAdminPassword == "" {
		c.iamClientErr = ErrMissingPassword
		return
	}
	err = client.Login(c.OrgAdminUsername, c.OrgAdminPassword)
	if err != nil {
		c.iamClientErr = err
		return
	}
	c.iamClient = client
	return
}

// setupS3CredsClient sets up an HSDP S3 Credentials client
func (c *Config) setupS3CredsClient() {
	if c.iamClientErr != nil {
		c.credsClient = nil
		c.credsClientErr = c.iamClientErr
		return
	}
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
