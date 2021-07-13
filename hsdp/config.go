package hsdp

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/fhir/go/jsonformat"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/philips-software/go-hsdp-api/cartel"
	"github.com/philips-software/go-hsdp-api/cdl"
	"github.com/philips-software/go-hsdp-api/cdr"
	"github.com/philips-software/go-hsdp-api/config"
	"github.com/philips-software/go-hsdp-api/console"
	"github.com/philips-software/go-hsdp-api/dicom"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/philips-software/go-hsdp-api/notification"
	"github.com/philips-software/go-hsdp-api/pki"
	"github.com/philips-software/go-hsdp-api/s3creds"
	"github.com/philips-software/go-hsdp-api/stl"
)

// Config contains configuration for the client
type Config struct {
	iam.Config
	BuildVersion      string
	ServiceID         string
	ServicePrivateKey string
	S3CredsURL        string
	NotificationURL   string
	STLURL            string
	CartelHost        string
	CartelToken       string
	CartelSecret      string
	CartelNoTLS       bool
	CartelSkipVerify  bool
	RetryMax          int
	UAAUsername       string
	UAAPassword       string
	UAAURL            string

	iamClient             *iam.Client
	cartelClient          *cartel.Client
	s3credsClient         *s3creds.Client
	consoleClient         *console.Client
	pkiClient             *pki.Client
	stlClient             *stl.Client
	notificationClient    *notification.Client
	debugFile             *os.File
	credsClientErr        error
	cartelClientErr       error
	iamClientErr          error
	consoleClientErr      error
	pkiClientErr          error
	stlClientErr          error
	notificationClientErr error
	TimeZone              string

	ma *jsonformat.Marshaller
	um *jsonformat.Unmarshaller
}

func (c *Config) IAMClient() (*iam.Client, error) {
	return c.iamClient, c.iamClientErr
}

func (c *Config) CartelClient() (*cartel.Client, error) {
	return c.cartelClient, c.cartelClientErr
}

func (c *Config) S3CredsClient() (*s3creds.Client, error) {
	return c.s3credsClient, c.credsClientErr
}

func (c *Config) ConsoleClient() (*console.Client, error) {
	return c.consoleClient, c.consoleClientErr
}

func (c *Config) STLClient(endpoint ...string) (*stl.Client, error) {
	return c.stlClient, c.stlClientErr
}

func (c *Config) PKIClient(regionEnvironment ...string) (*pki.Client, error) {
	if len(regionEnvironment) == 2 && c.consoleClient != nil && c.iamClient != nil {
		region := regionEnvironment[0]
		environment := regionEnvironment[1]
		return pki.NewClient(c.consoleClient, c.iamClient, &pki.Config{
			Region:      region,
			Environment: environment,
			DebugLog:    c.DebugLog,
		})
	}
	return c.pkiClient, c.pkiClientErr
}

func (c *Config) S3CredsClientWithLogin(username, password string) (*s3creds.Client, error) {
	if c.iamClientErr != nil {
		return nil, c.iamClientErr
	}
	newIAMClient, err := c.iamClient.WithLogin(username, password)
	if err != nil {
		return nil, err
	}
	return s3creds.NewClient(newIAMClient, &s3creds.Config{
		BaseURL:  c.S3CredsURL,
		DebugLog: c.DebugLog,
	})
}

func (c *Config) NotificationClient() (*notification.Client, error) {
	return c.notificationClient, c.notificationClientErr
}

// setupIAMClient sets up an HSDP IAM client
func (c *Config) setupIAMClient() {
	var standardClient *http.Client
	if c.RetryMax > 0 {
		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = c.RetryMax
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

func (c *Config) setupSTLClient() {
	if c.consoleClientErr != nil {
		c.stlClient = nil
		c.stlClientErr = c.consoleClientErr
		return
	}
	region := c.Region
	if region == "" {
		region = "dev"
	}
	ac, err := config.New(config.WithRegion(region))
	if err == nil {
		if url := ac.Service("stl").URL; c.STLURL == "" {
			c.STLURL = url
		}
	}
	client, err := stl.NewClient(c.consoleClient, &stl.Config{
		STLAPIURL: c.STLURL,
		DebugLog:  c.DebugLog,
	})
	if err != nil {
		c.stlClient = nil
		c.stlClientErr = err
		return
	}
	c.stlClient = client
}

func (c *Config) setupS3CredsClient() {
	if c.iamClientErr != nil {
		c.s3credsClient = nil
		c.credsClientErr = c.iamClientErr
		return
	}
	if c.Region != "" {
		env := c.Environment
		if env == "" {
			env = "prod"
		}
		ac, err := config.New(config.WithRegion(c.Region), config.WithEnv(env))
		if err == nil {
			if url := ac.Service("s3creds").URL; c.S3CredsURL == "" {
				c.S3CredsURL = url
			}
		}
	}
	client, err := s3creds.NewClient(c.iamClient, &s3creds.Config{
		BaseURL:  c.S3CredsURL,
		DebugLog: c.DebugLog,
	})
	if err != nil {
		c.s3credsClient = nil
		c.credsClientErr = err
		return
	}
	c.s3credsClient = client
}

func (c *Config) setupNotificationClient() {
	if c.iamClientErr != nil {
		c.notificationClient = nil
		c.notificationClientErr = c.iamClientErr
		return
	}
	if c.Region != "" {
		env := c.Environment
		if env == "" {
			env = "prod"
		}
		ac, err := config.New(config.WithRegion(c.Region), config.WithEnv(env))
		if err == nil {
			if url := ac.Service("notification").URL; c.NotificationURL == "" {
				c.NotificationURL = url
			}
		}
	}
	client, err := notification.NewClient(c.iamClient, &notification.Config{
		NotificationURL: c.NotificationURL,
		DebugLog:        c.DebugLog,
	})
	if err != nil {
		c.notificationClient = nil
		c.notificationClientErr = err
		return
	}
	c.notificationClient = client
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

// getCDLClient creates a HSDP CDL client
func (c *Config) getCDLClient(baseURL, tenantID string) (*cdl.Client, error) {
	if c.iamClientErr != nil {
		return nil, fmt.Errorf("IAM client error in getCDLClient: %w", c.iamClientErr)
	}
	if tenantID == "" {
		return nil, fmt.Errorf("getCDLClient: %w", ErrMissingOrganizationID)
	}
	client, err := cdl.NewClient(c.iamClient, &cdl.Config{
		CDLURL:         baseURL,
		OrganizationID: tenantID,
		DebugLog:       c.DebugLog,
	})
	if err != nil {
		return nil, fmt.Errorf("getFHIRClient: %w", err)
	}
	return client, nil
}

// getFHIRClient creates a HSDP CDR client
func (c *Config) getFHIRClient(baseURL, rootOrgID string) (*cdr.Client, error) {
	if c.iamClientErr != nil {
		return nil, fmt.Errorf("IAM client error in getFHIRClient: %w", c.iamClientErr)
	}
	if rootOrgID == "" {
		return nil, fmt.Errorf("getFHIRClient: %w", ErrMissingOrganizationID)
	}
	client, err := cdr.NewClient(c.iamClient, &cdr.Config{
		CDRURL:    baseURL,
		RootOrgID: rootOrgID,
		TimeZone:  c.TimeZone,
		DebugLog:  c.DebugLog,
	})
	if err != nil {
		return nil, fmt.Errorf("getFHIRClient: %w", err)
	}
	return client, nil
}

func (c *Config) Debug(format string, a ...interface{}) (int, error) {
	if c.debugFile != nil {
		output := fmt.Sprintf(format, a...)
		return c.debugFile.WriteString(output)
	}
	return 0, nil
}

func (c *Config) getDICOMConfigClient(url string) (*dicom.Client, error) {
	if c.iamClientErr != nil {
		return nil, fmt.Errorf("DICM client error in getDICOMConfigClient: %w", c.iamClientErr)
	}
	client, err := dicom.NewClient(c.iamClient, &dicom.Config{
		DICOMConfigURL: url,
		TimeZone:       c.TimeZone,
		DebugLog:       c.DebugLog,
	})
	if err != nil {
		return nil, fmt.Errorf("getDICOMConfigClient: %w", err)
	}
	return client, nil
}

func (c *Config) setupPKIClient() {
	if c.iamClientErr != nil {
		c.pkiClientErr = fmt.Errorf("IAM client error in setupPKIClient: %w", c.iamClientErr)
		return
	}
	if c.consoleClientErr != nil {
		c.pkiClientErr = fmt.Errorf("console client error in setupPKIClient: %w", c.consoleClientErr)
		return
	}
	client, err := pki.NewClient(c.consoleClient, c.iamClient, &pki.Config{
		Region:      c.Region,
		Environment: c.Environment,
		DebugLog:    c.DebugLog,
	})
	if err != nil {
		c.pkiClient = nil
		c.pkiClientErr = err
		return
	}
	c.pkiClient = client
	c.pkiClientErr = nil
}
