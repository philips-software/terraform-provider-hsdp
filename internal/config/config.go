package config

import (
	"fmt"
	"io"
	"net/http"

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
	"github.com/philips-software/go-hsdp-api/discovery"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/philips-software/go-hsdp-api/notification"
	"github.com/philips-software/go-hsdp-api/pki"
	"github.com/philips-software/go-hsdp-api/s3creds"
	"github.com/philips-software/go-hsdp-api/stl"
)

// Config contains configuration for the client
type Config struct {
	BuildVersion        string    `json:"-"`
	ServiceID           string    `json:"service_id"`
	ServicePrivateKey   string    `json:"service_private_key"`
	S3CredsURL          string    `json:"s3_creds_url"`
	NotificationURL     string    `json:"notification_url"`
	IAMURL              string    `json:"iam_url"`
	IDMURL              string    `json:"idm_url"`
	SharedKey           string    `json:"shared_key"`
	SecretKey           string    `json:"secret_key"`
	MDMURL              string    `json:"mdm_url"`
	Region              string    `json:"region"`
	Environment         string    `json:"environment"`
	OAuth2ClientID      string    `json:"oauth2_client_id"`
	OAuth2ClientSecret  string    `json:"oauth2_client_secret"`
	STLURL              string    `json:"stl_url"`
	OrgAdminUsername    string    `json:"org_admin_username"`
	OrgAdminPassword    string    `json:"org_admin_password"`
	DebugLog            string    `json:"debug_log"`
	DebugWriter         io.Writer `json:"-"`
	CartelHost          string    `json:"cartel_host"`
	CartelToken         string    `json:"cartel_token"`
	CartelSecret        string    `json:"cartel_secret"`
	CartelNoTLS         bool      `json:"cartel_no_tls"`
	CartelSkipVerify    bool      `json:"cartel_skip_verify"`
	RetryMax            int       `json:"retry_max"`
	UAAUsername         string    `json:"uaa_username"`
	UAAPassword         string    `json:"uaa_password"`
	UAAURL              string    `json:"uaa_url"`
	AIInferenceEndpoint string    `json:"ai_inference_endpoint"`
	AIWorkspaceEndpoint string    `json:"ai_workspace_endpoint"`

	iamClient             *iam.Client
	cartelClient          *cartel.Client
	s3credsClient         *s3creds.Client
	consoleClient         *console.Client
	pkiClient             *pki.Client
	stlClient             *stl.Client
	notificationClient    *notification.Client
	mdmClient             *mdm.Client
	discoveryClient       *discovery.Client
	DebugStdErr           bool `json:"debugging"`
	credsClientErr        error
	cartelClientErr       error
	iamClientErr          error
	consoleClientErr      error
	pkiClientErr          error
	stlClientErr          error
	notificationClientErr error
	mdmClientErr          error
	discoveryClientErr    error
	TimeZone              string `json:"time_zone"`

	STU3MA *jsonformat.Marshaller   `json:"-"`
	STU3UM *jsonformat.Unmarshaller `json:"-"`
	R4MA   *jsonformat.Marshaller   `json:"-"`
	R4UM   *jsonformat.Unmarshaller `json:"-"`
}

func (c *Config) IAMClient(principal ...*Principal) (*iam.Client, error) {
	if len(principal) > 0 && principal[0] != nil && principal[0].HasAuth() {
		p := principal[0]
		cfg := iam.Config{
			OAuth2ClientID: c.OAuth2ClientID,
			OAuth2Secret:   c.OAuth2ClientSecret,
			Region:         c.Region,
			Environment:    c.Environment,
			DebugLog:       c.DebugWriter,
			SharedKey:      c.SharedKey,
			SecretKey:      c.SecretKey,
			IDMURL:         c.IDMURL,
			IAMURL:         c.IAMURL,
		}
		if p.OAuth2ClientID != "" {
			cfg.OAuth2ClientID = p.OAuth2ClientID
		}
		if p.OAuth2Password != "" {
			cfg.OAuth2Secret = p.OAuth2Password
		}
		if p.Environment != "" {
			cfg.Environment = p.Environment
		}
		if p.Region != "" {
			cfg.Region = p.Region
		}
		iamClient, err := iam.NewClient(nil, &cfg)
		if err != nil {
			return nil, err
		}
		if p.Username != "" {
			err := iamClient.Login(p.Username, p.Password)
			if err != nil {
				return nil, err
			}
			return iamClient, nil
		}
		if p.ServiceID != "" {
			err := iamClient.ServiceLogin(iam.Service{
				ServiceID:  p.ServiceID,
				PrivateKey: p.ServicePrivateKey,
			})
			if err != nil {
				return nil, err
			}
		}
		return iamClient, nil
	}
	return c.iamClient, c.iamClientErr
}

func (c *Config) HasUAAuth() bool {
	return c.UAAUsername != "" && c.UAAPassword != ""
}

func (c *Config) DiscoveryClient(principal ...*Principal) (*discovery.Client, error) {
	if len(principal) > 0 && principal[0] != nil && principal[0].HasAuth() {
		region := principal[0].Region
		environment := principal[0].Environment
		iamClient, err := c.IAMClient(principal...)
		if err != nil {
			return nil, err
		}
		return discovery.NewClient(iamClient, &discovery.Config{
			Region:      region,
			Environment: environment,
			DebugLog:    c.DebugWriter,
		})
	}
	return c.discoveryClient, c.discoveryClientErr
}

func (c *Config) CartelClient() (*cartel.Client, error) {
	return c.cartelClient, c.cartelClientErr
}

func (c *Config) S3CredsClient() (*s3creds.Client, error) {
	return c.s3credsClient, c.credsClientErr
}

func (c *Config) ConsoleClient(principal ...*Principal) (*console.Client, error) {
	region := c.Region
	uaaUsername := c.UAAUsername
	uaaPassword := c.UAAPassword

	if region == "" {
		region = "dev"
	}

	if len(principal) == 0 && c.consoleClient != nil {
		return c.consoleClient, c.consoleClientErr
	}

	if len(principal) > 0 && principal[0] != nil {
		p := principal[0]
		if p.Region != "" {
			region = p.Region
		}
		if p.UAAUsername != "" {
			uaaUsername = p.UAAUsername
		}
		if p.UAAPassword != "" {
			uaaPassword = p.UAAPassword
		}
	}
	client, err := console.NewClient(nil, &console.Config{
		Region:   region,
		DebugLog: c.DebugWriter,
	})

	if err != nil {
		return nil, err
	}
	if uaaUsername == "" || uaaPassword == "" {
		return nil, ErrMissingUAACredentials
	}
	err = client.Login(uaaUsername, uaaPassword)
	if err != nil {
		return nil, err
	}
	return client, err
}

func (c *Config) MDMClient() (*mdm.Client, error) {
	return c.mdmClient, c.mdmClientErr
}

func (c *Config) STLClient(principal ...*Principal) (*stl.Client, error) {
	region := c.Region
	if region == "" {
		region = "dev"
	}
	consoleClient := c.consoleClient
	consoleClientErr := c.consoleClientErr
	stlURL := c.STLURL

	if len(principal) == 0 {
		return c.stlClient, c.stlClientErr
	}

	if principal[0] != nil {
		p := principal[0]
		if p.Region != "" {
			region = p.Region
			ac, err := config.New(config.WithRegion(region))
			if err == nil {
				if url := ac.Service("stl").URL; url != "" {
					stlURL = url
				}
			}
		}
		if p.Endpoint != "" {
			stlURL = p.Endpoint
		}
		consoleClient, consoleClientErr = c.ConsoleClient(principal...)
	}
	if consoleClientErr != nil {
		return nil, consoleClientErr
	}

	client, err := stl.NewClient(consoleClient, &stl.Config{
		Region:    region,
		STLAPIURL: stlURL,
		DebugLog:  c.DebugWriter,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Config) DockerClient(principal ...*Principal) (*docker.Client, error) {
	r := c.Region
	if len(principal) > 0 && principal[0] != nil {
		r = principal[0].Region
	}
	if c.consoleClientErr != nil {
		return nil, c.consoleClientErr
	}
	return docker.NewClient(c.consoleClient, &docker.Config{
		Region: r,
	})
}

func (c *Config) PKIClient(principal ...*Principal) (*pki.Client, error) {
	if len(principal) > 0 && principal[0] != nil && principal[0].HasAuth() && c.consoleClient != nil {
		region := principal[0].Region
		environment := principal[0].Environment
		iamClient, err := c.IAMClient(principal...)
		if err != nil {
			return nil, err
		}
		return pki.NewClient(c.consoleClient, iamClient, &pki.Config{
			Region:      region,
			Environment: environment,
			DebugLog:    c.DebugWriter,
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
		DebugLog: c.DebugWriter,
	})
}

func (c *Config) NotificationClient(principal ...*Principal) (*notification.Client, error) {
	if len(principal) > 0 && principal[0] != nil && principal[0].HasAuth() {
		region := principal[0].Region
		environment := principal[0].Environment
		iamClient, err := c.IAMClient(principal...)
		if err != nil {
			return nil, err
		}
		endpoint := principal[0].Endpoint
		if endpoint == "" {
			ac, err := config.New(config.WithRegion(region), config.WithEnv(environment))
			if err == nil {
				if url := ac.Service("notification").URL; url != "" {
					endpoint = url
				}
			}
		}
		return notification.NewClient(iamClient, &notification.Config{
			Region:          region,
			Environment:     environment,
			NotificationURL: endpoint,
			DebugLog:        c.DebugWriter,
		})
	}
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
	cfg := &iam.Config{
		OAuth2ClientID: c.OAuth2ClientID,
		OAuth2Secret:   c.OAuth2ClientSecret,
		Region:         c.Region,
		Environment:    c.Environment,
		DebugLog:       c.DebugWriter,
		SharedKey:      c.SharedKey,
		SecretKey:      c.SecretKey,
		IDMURL:         c.IDMURL,
		IAMURL:         c.IAMURL,
	}
	client, err := iam.NewClient(standardClient, cfg)
	if err != nil {
		c.iamClientErr = fmt.Errorf("possible invalid environment/region: %w", err)
		return
	}
	usingServiceIdentity := false
	if c.ServiceID != "" && c.ServicePrivateKey != "" {
		err = client.ServiceLogin(iam.Service{
			ServiceID:  c.ServiceID,
			PrivateKey: c.ServicePrivateKey,
		})
		if err != nil {
			c.iamClientErr = fmt.Errorf("invalid IAM Service Identity credentials: %w", err)
			return
		}
		usingServiceIdentity = true
	}
	usingOrgAdmin := false
	if !usingServiceIdentity && c.OrgAdminUsername != "" && c.OrgAdminPassword != "" {
		if c.OAuth2ClientID == "" {
			c.iamClientErr = ErrMissingClientID
			return
		}
		err = client.Login(c.OrgAdminUsername, c.OrgAdminPassword)
		if err != nil {
			c.iamClientErr = fmt.Errorf("invalid IAM Org Admin credentials: %w", err)
			return
		}
		usingOrgAdmin = true
	}
	if !(usingServiceIdentity || usingOrgAdmin) {
		c.iamClientErr = fmt.Errorf("invalid / missing IAM Service Identity or IAM Org Admin credentials")
		return
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
		DebugLog:  c.DebugWriter,
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
		DebugLog: c.DebugWriter,
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
		DebugLog:        c.DebugWriter,
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
			url := ac.Service("connect-mdm").URL
			if url != "" {
				c.MDMURL = url
			} else {
				c.mdmClient = nil
				c.mdmClientErr = fmt.Errorf("missing MDM URL (%s/%s), you can set a custom value using 'mdm_url'", env, c.Region)
				return
			}
		}
	}
	client, err := mdm.NewClient(c.iamClient, &mdm.Config{
		BaseURL:  c.MDMURL,
		DebugLog: c.DebugWriter,
	})
	if err != nil {
		c.mdmClient = nil
		c.mdmClientErr = fmt.Errorf("configuration error (%s/%s): %w", c.Environment, c.Region, err)
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
	if c.CartelToken == "" || c.CartelSecret == "" {
		c.cartelClient = nil
		c.cartelClientErr = fmt.Errorf("missing Cartel token or secret, set 'cartel_token' and 'cartel_secret'")
		return
	}
	client, err := cartel.NewClient(nil, &cartel.Config{
		Region:     c.Region,
		Host:       c.CartelHost,
		Token:      c.CartelToken,
		Secret:     c.CartelSecret,
		NoTLS:      c.CartelNoTLS,
		SkipVerify: c.CartelSkipVerify,
		DebugLog:   c.DebugWriter,
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
		DebugLog: c.DebugWriter,
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
		c.consoleClient = nil
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
		DebugLog:  c.DebugWriter,
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
		DebugLog: c.DebugWriter,
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
		DebugLog:       c.DebugWriter,
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
		DebugLog:       c.DebugWriter,
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
		DebugLog:       c.DebugWriter,
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
		DebugLog:       c.DebugWriter,
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
		DebugLog:       c.DebugWriter,
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
		DebugLog:  c.DebugWriter,
	})
	if err != nil {
		return nil, fmt.Errorf("GetFHIRClient: %w", err)
	}
	return client, nil
}

func (c *Config) Debug(format string, a ...interface{}) (int, error) {
	if c.DebugWriter != nil {
		output := fmt.Sprintf(format, a...)
		return io.WriteString(c.DebugWriter, output)
	}
	return 0, nil
}

func (c *Config) GetDICOMConfigClient(url string) (*dicom.Client, error) {
	if c.iamClientErr != nil {
		return nil, fmt.Errorf("DICM client error in GetDICOMConfigClient: %w", c.iamClientErr)
	}
	if url == "" {
		return nil, fmt.Errorf("GetDICOMConfigClient: empty config_url")
	}
	client, err := dicom.NewClient(c.iamClient, &dicom.Config{
		DICOMConfigURL: url,
		TimeZone:       c.TimeZone,
		DebugLog:       c.DebugWriter,
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
	// We ignore any consoleClient error for now
	client, err := pki.NewClient(c.consoleClient, c.iamClient, &pki.Config{
		Region:      c.Region,
		Environment: c.Environment,
		DebugLog:    c.DebugWriter,
	})
	if err != nil {
		c.pkiClient = nil
		c.pkiClientErr = err
		return
	}
	c.pkiClient = client
	c.pkiClientErr = nil
}

func (c *Config) SetupDiscoveryClient() {
	if c.iamClientErr != nil {
		c.pkiClientErr = fmt.Errorf("IAM client error in SetupDiscoveryClient: %w", c.iamClientErr)
		return
	}
	client, err := discovery.NewClient(c.iamClient, &discovery.Config{
		Region:      c.Region,
		Environment: c.Environment,
		DebugLog:    c.DebugWriter,
	})
	if err != nil {
		c.discoveryClient = nil
		c.discoveryClientErr = err
		return
	}
	c.discoveryClient = client
	c.discoveryClientErr = nil
}
