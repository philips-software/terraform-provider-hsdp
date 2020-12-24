package hsdp

import (
	"github.com/google/fhir/go/jsonformat"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/philips-software/go-hsdp-api/cartel"
	"github.com/philips-software/go-hsdp-api/cdr"
	"github.com/philips-software/go-hsdp-api/config"
	"github.com/philips-software/go-hsdp-api/console"
	"github.com/philips-software/go-hsdp-api/credentials"
	"github.com/philips-software/go-hsdp-api/iam"
	"net/http"
)

// Config contains configuration for the client
type Config struct {
	iam.Config
	ServiceID         string
	ServicePrivateKey string
	S3CredsURL        string
	CartelHost        string
	CartelToken       string
	CartelSecret      string
	CartelNoTLS       bool
	CartelSkipVerify  bool
	RetryMax          int
	UAAUsername       string
	UAAPassword       string
	UAAURL            string

	iamClient        *iam.Client
	cartelClient     *cartel.Client
	credsClient      *credentials.Client
	consoleClient    *console.Client
	credsClientErr   error
	cartelClientErr  error
	iamClientErr     error
	consoleClientErr error
	TimeZone         string

	ma *jsonformat.Marshaller
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

func (c *Config) ConsoleClient() (*console.Client, error) {
	return c.consoleClient, c.consoleClientErr
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
	if c.ServiceID != "" && c.ServicePrivateKey != "" {
		err = client.ServiceLogin(iam.Service{
			ServiceID:  c.ServiceID,
			PrivateKey: c.ServicePrivateKey,
		})
		if err != nil {
			c.iamClientErr = err
			return
		}
	}
	if c.OrgAdminUsername != "" && c.OrgAdminPassword != "" {
		err = client.Login(c.OrgAdminUsername, c.OrgAdminPassword)
		if err != nil {
			c.iamClientErr = err
			return
		}
	}
	c.iamClient = client
}

// setupS3CredsClient sets up an HSDP S3 Credentials client
func (c *Config) setupS3CredsClient() {
	if c.iamClientErr != nil {
		c.credsClient = nil
		c.credsClientErr = c.iamClientErr
		return
	}
	if c.Environment != "" && c.Region != "" {
		ac, err := config.New(config.WithRegion(c.Region), config.WithEnv(c.Environment))
		if err == nil {
			if url := ac.Service("s3creds").URL; c.S3CredsURL == "" {
				c.S3CredsURL = url
			}
		}
	}
	client, err := credentials.NewClient(c.iamClient, &credentials.Config{
		BaseURL:  c.S3CredsURL,
		DebugLog: c.DebugLog,
	})
	if err != nil {
		c.credsClient = nil
		c.credsClientErr = err
		return
	}
	c.credsClient = client
}

// setupCartelClient sets up an Cartel client
func (c *Config) setupCartelClient() {
	client, err := cartel.NewClient(nil, &cartel.Config{
		Region:     c.Region,
		Host:       c.CartelHost,
		Token:      c.CartelToken,
		Secret:     c.CartelSecret,
		NoTLS:      c.CartelNoTLS,
		SkipVerify: c.CartelSkipVerify,
		DebugLog:   c.DebugLog,
	})
	if err != nil {
		c.cartelClient = nil
		c.cartelClientErr = err
		return
	}
	c.cartelClient = client
}

// setupConsoleClient sets up an Console client
func (c *Config) setupConsoleClient() {
	client, err := console.NewClient(nil, &console.Config{
		Region:   c.Region,
		DebugLog: c.DebugLog,
	})
	if err != nil {
		c.consoleClient = nil
		c.consoleClientErr = err
		return
	}
	if c.UAAUsername != "" && c.UAAPassword != "" {
		err = client.Login(c.UAAUsername, c.UAAPassword)
		if err != nil {
			c.consoleClientErr = err
			return
		}
	}
	c.consoleClient = client
}

// getFHIRClientFromEndpoint creates a HSDP CDR client form the given endpoint
func (c *Config) getFHIRClientFromEndpoint(endpointURL string) (*cdr.Client, error) {
	if c.iamClientErr != nil {
		return nil, c.iamClientErr
	}
	client, err := cdr.NewClient(c.iamClient, &cdr.Config{
		CDRURL:    "https://localhost.domain",
		RootOrgID: "",
		TimeZone:  c.TimeZone,
		DebugLog:  c.DebugLog,
	})
	if err != nil {
		return nil, err
	}
	if err = client.SetEndpointURL(endpointURL); err != nil {
		return nil, err
	}
	return client, nil
}

// getFHIRClient creates a HSDP CDR client
func (c *Config) getFHIRClient(baseURL, rootOrgID string) (*cdr.Client, error) {
	if c.iamClientErr != nil {
		return nil, c.iamClientErr
	}
	if rootOrgID == "" {
		return nil, ErrMissingOrganizationID
	}
	client, err := cdr.NewClient(c.iamClient, &cdr.Config{
		CDRURL:    baseURL,
		RootOrgID: rootOrgID,
		TimeZone:  c.TimeZone,
		DebugLog:  c.DebugLog,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
