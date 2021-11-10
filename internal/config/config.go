package config

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/fhir/go/jsonformat"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/philips-software/go-hsdp-api/ai"
	"github.com/philips-software/go-hsdp-api/ai/inference"
	"github.com/philips-software/go-hsdp-api/ai/workspace"
	"github.com/philips-software/go-hsdp-api/cartel"
	"github.com/philips-software/go-hsdp-api/cdl"
	"github.com/philips-software/go-hsdp-api/cdr"
	"github.com/philips-software/go-hsdp-api/config"
	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/philips-software/go-hsdp-api/console"
	"github.com/philips-software/go-hsdp-api/console/docker"
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
	BuildVersion        string
	ServiceID           string
	ServicePrivateKey   string
	S3CredsURL          string
	NotificationURL     string
	MDMURL              string
	STLURL              string
	CartelHost          string
	CartelToken         string
	CartelSecret        string
	CartelNoTLS         bool
	CartelSkipVerify    bool
	RetryMax            int
	UAAUsername         string
	UAAPassword         string
	UAAURL              string
	AIInferenceEndpoint string
	AIWorkspaceEndpoint string

	iamClient             *iam.Client
	cartelClient          *cartel.Client
	s3credsClient         *s3creds.Client
	consoleClient         *console.Client
	pkiClient             *pki.Client
	stlClient             *stl.Client
	notificationClient    *notification.Client
	mdmClient             *mdm.Client
	DebugFile             *os.File
	credsClientErr        error
	cartelClientErr       error
	iamClientErr          error
	consoleClientErr      error
	pkiClientErr          error
	stlClientErr          error
	notificationClientErr error
	mdmClientErr          error
	TimeZone              string

	STU3MA *jsonformat.Marshaller
	STU3UM *jsonformat.Unmarshaller
	R4MA   *jsonformat.Marshaller
	R4UM   *jsonformat.Unmarshaller
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

func (c *Config) MDMClient() (*mdm.Client, error) {
	return c.mdmClient, c.mdmClientErr
}

func (c *Config) STLClient(_ ...string) (*stl.Client, error) {
	return c.stlClient, c.stlClientErr
}

func (c *Config) DockerClient(region ...string) (*docker.Client, error) {
	r := c.Region
	if len(region) > 0 {
		r = region[0]
	}
	if c.consoleClientErr != nil {
		return nil, c.consoleClientErr
	}
	return docker.NewClient(c.consoleClient, &docker.Config{
		Region: r,
	})
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

// SetupIAMClient sets up an HSDP IAM client
func (c *Config) SetupIAMClient() {
	var standardClient *http.Client
	if c.RetryMax > 0 {
		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = c.RetryMax
		standardClient = retryClient.StandardClient()
	}
	c.iamClient = nil
	client, err := iam.NewClient(standardClient, &c.Config)
	if err != nil {
		c.iamClientErr = fmt.Errorf("possible invalid region/environment: %w", err)
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
		if c.OAuth2ClientID == "" {
			c.iamClientErr = ErrMissingClientID
			return
		}
		err = client.Login(c.OrgAdminUsername, c.OrgAdminPassword)
		if err != nil {
			c.iamClientErr = err
			return
		}
	}
	c.iamClient = client
}

func (c *Config) SetupSTLClient() {
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

func (c *Config) SetupS3CredsClient() {
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

func (c *Config) SetupNotificationClient() {
	if c.iamClientErr != nil {
		c.notificationClient = nil
		c.notificationClientErr = c.iamClientErr
		return
	}
	if c.NotificationURL == "" {
		env := c.Environment
		if env == "" {
			env = "prod"
		}
		ac, err := config.New(config.WithRegion(c.Region), config.WithEnv(env))
		if err == nil {
			if url := ac.Service("notification").URL; url != "" {
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

func (c *Config) SetupMDMClient() {
	if c.iamClientErr != nil {
		c.mdmClient = nil
		c.mdmClientErr = c.iamClientErr
		return
	}
	if c.MDMURL == "" {
		env := c.Environment
		if env == "" {
			env = "prod"
		}
		ac, err := config.New(config.WithRegion(c.Region), config.WithEnv(env))
		if err == nil {
			if url := ac.Service("connect-mdm").URL; url != "" {
				c.MDMURL = url
			}
		}
	}
	client, err := mdm.NewClient(c.iamClient, &mdm.Config{
		BaseURL:  c.MDMURL,
		DebugLog: c.DebugLog,
	})
	if err != nil {
		c.mdmClient = nil
		c.mdmClientErr = err
		return
	}
	c.mdmClient = client
}

// SetupCartelClient sets up an Cartel client
func (c *Config) SetupCartelClient() {
	if c.CartelHost == "" {
		ac, err := config.New(config.WithRegion(c.Region))
		if err == nil {
			if host := ac.Service("cartel").Host; host != "" {
				c.CartelHost = host
			}
		}
	}
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

// SetupConsoleClient sets up an Console client
func (c *Config) SetupConsoleClient() {
	client, err := console.NewClient(nil, &console.Config{
		Region:   c.Region,
		DebugLog: c.DebugLog,
	})
	if err != nil {
		c.consoleClient = nil
		c.consoleClientErr = err
		return
	}
	if c.UAAUsername == "" || c.UAAPassword == "" {
		c.consoleClientErr = ErrMissingUAACredentials
		c.consoleClient = nil
		return
	}
	err = client.Login(c.UAAUsername, c.UAAPassword)
	if err != nil {
		c.consoleClientErr = err
		return
	}
	c.consoleClient = client
}

func (c *Config) GetFHIRClientFromEndpoint(endpointURL string) (*cdr.Client, error) {
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

func (c *Config) GetCDLClientFromEndpoint(endpointURL string) (*cdl.Client, error) {
	if c.iamClientErr != nil {
		return nil, c.iamClientErr
	}
	client, err := cdl.NewClient(c.iamClient, &cdl.Config{
		CDLURL:   "https://localhost.domain",
		DebugLog: c.DebugLog,
	})
	if err != nil {
		return nil, err
	}
	if err = client.SetEndpointURL(endpointURL); err != nil {
		return nil, err
	}
	return client, nil
}

// GetCDLClient creates a HSDP CDL client
func (c *Config) GetCDLClient(baseURL, tenantID string) (*cdl.Client, error) {
	if c.iamClientErr != nil {
		return nil, fmt.Errorf("IAM client error in GetCDLClient: %w", c.iamClientErr)
	}
	if tenantID == "" {
		return nil, fmt.Errorf("GetCDLClient: %w", ErrMissingOrganizationID)
	}
	client, err := cdl.NewClient(c.iamClient, &cdl.Config{
		CDLURL:         baseURL,
		OrganizationID: tenantID,
		DebugLog:       c.DebugLog,
	})
	if err != nil {
		return nil, fmt.Errorf("GetCDLClient: %w", err)
	}
	return client, nil
}

func (c *Config) GetAIInferenceClient(baseURL, tenantID string) (*inference.Client, error) {
	if c.iamClientErr != nil {
		return nil, fmt.Errorf("IAM client error in getAIInferenceClient: %w", c.iamClientErr)
	}
	if tenantID == "" {
		return nil, fmt.Errorf("getAIInferenceClient: %w", ErrMissingOrganizationID)
	}
	client, err := inference.NewClient(c.iamClient, &ai.Config{
		BaseURL:        baseURL,
		OrganizationID: tenantID,
		DebugLog:       c.DebugLog,
	})
	if err != nil {
		return nil, fmt.Errorf("getAIInferenceClient: %w", err)
	}
	return client, nil
}

func (c *Config) GetAIInferenceClientFromEndpoint(endpointURL string) (*inference.Client, error) {
	if c.iamClientErr != nil {
		return nil, c.iamClientErr
	}
	if endpointURL == "" {
		endpointURL = c.AIInferenceEndpoint
	}
	client, err := inference.NewClient(c.iamClient, &ai.Config{
		BaseURL:        "http://localhost",
		OrganizationID: "not-set",
		DebugLog:       c.DebugLog,
	})
	if err != nil {
		return nil, fmt.Errorf("getAIInferenceClientFromEndpoint: %w", err)
	}
	if err = client.SetEndpointURL(endpointURL); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Config) GetAIWorkspaceClient(baseURL, tenantID string) (*workspace.Client, error) {
	if c.iamClientErr != nil {
		return nil, fmt.Errorf("IAM client error in getAIWorkspaceClient: %w", c.iamClientErr)
	}
	if tenantID == "" {
		return nil, fmt.Errorf("getAIWorkspaceClient: %w", ErrMissingOrganizationID)
	}
	client, err := workspace.NewClient(c.iamClient, &ai.Config{
		BaseURL:        baseURL,
		OrganizationID: tenantID,
		DebugLog:       c.DebugLog,
	})
	if err != nil {
		return nil, fmt.Errorf("getAIWorkspaceClient: %w", err)
	}
	return client, nil
}

func (c *Config) GetAIWorkspaceClientFromEndpoint(endpointURL string) (*workspace.Client, error) {
	if c.iamClientErr != nil {
		return nil, c.iamClientErr
	}
	if endpointURL == "" {
		endpointURL = c.AIWorkspaceEndpoint
	}
	client, err := workspace.NewClient(c.iamClient, &ai.Config{
		BaseURL:        "http://localhost",
		OrganizationID: "not-set",
		DebugLog:       c.DebugLog,
	})
	if err != nil {
		return nil, fmt.Errorf("getAIWorkspaceClientFromEndpoint: %w", err)
	}
	if err = client.SetEndpointURL(endpointURL); err != nil {
		return nil, err
	}
	return client, nil
}

// GetFHIRClient creates a HSDP CDR client
func (c *Config) GetFHIRClient(baseURL, rootOrgID string) (*cdr.Client, error) {
	if c.iamClientErr != nil {
		return nil, fmt.Errorf("IAM client error in GetFHIRClient: %w", c.iamClientErr)
	}
	if rootOrgID == "" {
		return nil, fmt.Errorf("GetFHIRClient: %w", ErrMissingOrganizationID)
	}
	client, err := cdr.NewClient(c.iamClient, &cdr.Config{
		CDRURL:    baseURL,
		RootOrgID: rootOrgID,
		TimeZone:  c.TimeZone,
		DebugLog:  c.DebugLog,
	})
	if err != nil {
		return nil, fmt.Errorf("GetFHIRClient: %w", err)
	}
	return client, nil
}

func (c *Config) Debug(format string, a ...interface{}) (int, error) {
	if c.DebugFile != nil {
		output := fmt.Sprintf(format, a...)
		return c.DebugFile.WriteString(output)
	}
	return 0, nil
}

func (c *Config) GetDICOMConfigClient(url string) (*dicom.Client, error) {
	if c.iamClientErr != nil {
		return nil, fmt.Errorf("DICM client error in GetDICOMConfigClient: %w", c.iamClientErr)
	}
	client, err := dicom.NewClient(c.iamClient, &dicom.Config{
		DICOMConfigURL: url,
		TimeZone:       c.TimeZone,
		DebugLog:       c.DebugLog,
	})
	if err != nil {
		return nil, fmt.Errorf("GetDICOMConfigClient: %w", err)
	}
	return client, nil
}

func (c *Config) SetupPKIClient() {
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
