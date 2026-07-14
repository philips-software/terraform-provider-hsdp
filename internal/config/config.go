package config

import (
	"fmt"
	"io"
	"net/http"

	"github.com/philips-software/go-dip-api/connect/dbs"
	"github.com/philips-software/go-dip-api/connect/provisioning"

	"github.com/philips-software/go-dip-api/connect/blr"

	"github.com/google/fhir/go/jsonformat"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/philips-software/go-dip-api/config"
	"github.com/philips-software/go-dip-api/connect/mdm"
	"github.com/philips-software/go-dip-api/discovery"
	"github.com/philips-software/go-dip-api/iam"
	"github.com/philips-software/go-dip-api/stl"
	"golang.org/x/oauth2"
)

// Config contains configuration for the client
type Config struct {
	BuildVersion       string    `json:"-"`
	ServiceID          string    `json:"service_id"`
	ServicePrivateKey  string    `json:"service_private_key"`
	IAMURL             string    `json:"iam_url"`
	IDMURL             string    `json:"idm_url"`
	SharedKey          string    `json:"shared_key"`
	SecretKey          string    `json:"secret_key"`
	MDMURL             string    `json:"mdm_url"`
	Region             string    `json:"region"`
	Environment        string    `json:"environment"`
	OAuth2ClientID     string    `json:"oauth2_client_id"`
	OAuth2ClientSecret string    `json:"oauth2_client_secret"`
	STLURL             string    `json:"stl_url"`
	OrgAdminUsername   string    `json:"org_admin_username"`
	OrgAdminPassword   string    `json:"org_admin_password"`
	DebugLog           string    `json:"debug_log"`
	DebugWriter        io.Writer `json:"-"`
	RetryMax           int       `json:"retry_max"`
	UAAUsername        string    `json:"uaa_username"`
	UAAPassword        string    `json:"uaa_password"`
	UAAURL             string    `json:"uaa_url"`

	iamClient             *iam.Client
	stlClient             *stl.Client
	blrClient             *blr.Client
	mdmClient             *mdm.Client
	discoveryClient       *discovery.Client
	dbsClient             *dbs.Client
	provisioningClient    *provisioning.Client
	DebugStdErr           bool `json:"debugging"`
	iamClientErr          error
	stlClientErr          error
	mdmClientErr          error
	discoveryClientErr    error
	blrClientErr          error
	dbsClientErr          error
	provisioningClientErr error
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

func (c *Config) BLRClient(principal ...*Principal) (*blr.Client, error) {
	if len(principal) > 0 && principal[0] != nil && principal[0].HasAuth() {
		region := principal[0].Region
		environment := principal[0].Environment
		iamClient, err := c.IAMClient(principal...)
		if err != nil {
			return nil, err
		}
		return blr.NewClient(iamClient, &blr.Config{
			Region:      region,
			Environment: environment,
			DebugLog:    c.DebugWriter,
		})
	}
	return c.blrClient, c.blrClientErr
}

type iamTokenSource struct {
	client *iam.Client
}

func (t *iamTokenSource) Token() (*oauth2.Token, error) {
	tokenStr, err := t.client.Token()
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken: tokenStr,
	}, nil
}

func (c *Config) MDMClient() (*mdm.Client, error) {
	return c.mdmClient, c.mdmClientErr
}

func (c *Config) STLClient(principal ...*Principal) (*stl.Client, error) {
	region := c.Region
	if region == "" {
		region = "dev"
	}
	stlURL := c.STLURL

	if len(principal) == 0 {
		return c.stlClient, c.stlClientErr
	}

	var iamClient *iam.Client
	var iamClientErr error

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
		iamClient, iamClientErr = c.IAMClient(principal...)
	} else {
		iamClient, iamClientErr = c.IAMClient()
	}
	if iamClientErr != nil {
		return nil, iamClientErr
	}
	if iamClient == nil {
		return nil, fmt.Errorf("IAM client not initialized")
	}

	client, err := stl.NewClient(&iamTokenSource{client: iamClient}, &stl.Config{
		Region:    region,
		STLAPIURL: stlURL,
		DebugLog:  c.DebugWriter,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Config) DBSClient(principal ...*Principal) (*dbs.Client, error) {
	if len(principal) > 0 && principal[0] != nil && principal[0].HasAuth() {
		region := principal[0].Region
		environment := principal[0].Environment
		iamClient, err := c.IAMClient(principal...)
		if err != nil {
			return nil, err
		}
		return dbs.NewClient(iamClient, &dbs.Config{
			Region:      region,
			Environment: environment,
			DebugLog:    c.DebugWriter,
		})
	}
	return c.dbsClient, c.dbsClientErr
}

func (c *Config) ProvisioningClient(principal ...*Principal) (*provisioning.Client, error) {
	if len(principal) > 0 && principal[0] != nil && principal[0].HasAuth() {
		region := principal[0].Region
		environment := principal[0].Environment
		iamClient, err := c.IAMClient(principal...)
		if err != nil {
			return nil, err
		}
		return provisioning.NewClient(iamClient, &provisioning.Config{
			Region:      region,
			Environment: environment,
			DebugLog:    c.DebugWriter,
		})
	}
	return c.provisioningClient, c.provisioningClientErr
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
			c.iamClientErr = fmt.Errorf("invalid IAM Service Identity credentials for '%s': %w", c.ServiceID, err)
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
			c.iamClientErr = fmt.Errorf("invalid IAM Org Admin credentials for '%s': %w", c.OrgAdminUsername, err)
			return
		}
		usingOrgAdmin = true
	}
	if !usingServiceIdentity && !usingOrgAdmin {
		c.iamClientErr = fmt.Errorf("invalid / missing IAM Service Identity or IAM Org Admin credentials")
		return
	}
	c.iamClient = client
}

func (c *Config) SetupSTLClient() {
	if c.iamClientErr != nil {
		c.stlClient = nil
		c.stlClientErr = c.iamClientErr
		return
	}
	if c.iamClient == nil {
		c.stlClient = nil
		c.stlClientErr = fmt.Errorf("IAM client not initialized")
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
	client, err := stl.NewClient(&iamTokenSource{client: c.iamClient}, &stl.Config{
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

func (c *Config) Debug(format string, a ...interface{}) (int, error) {
	if c.DebugWriter != nil {
		output := fmt.Sprintf(format, a...)
		return io.WriteString(c.DebugWriter, output)
	}
	return 0, nil
}

func (c *Config) SetupDiscoveryClient() {
	if c.iamClientErr != nil {
		c.discoveryClientErr = fmt.Errorf("IAM client error in SetupDiscoveryClient: %w", c.iamClientErr)
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

func (c *Config) SetupBLRClient() {
	if c.iamClientErr != nil {
		c.blrClientErr = fmt.Errorf("IAM client error in SetupBLRClient: %w", c.iamClientErr)
		return
	}
	client, err := blr.NewClient(c.iamClient, &blr.Config{
		Region:      c.Region,
		Environment: c.Environment,
		DebugLog:    c.DebugWriter,
	})
	if err != nil {
		c.blrClient = nil
		c.blrClientErr = err
		return
	}
	c.blrClient = client
	c.blrClientErr = nil
}

func (c *Config) SetupDBSClient() {
	if c.iamClientErr != nil {
		c.blrClientErr = fmt.Errorf("IAM client error in SetupDBSClient: %w", c.iamClientErr)
		return
	}
	client, err := dbs.NewClient(c.iamClient, &dbs.Config{
		Region:      c.Region,
		Environment: c.Environment,
		DebugLog:    c.DebugWriter,
	})
	if err != nil {
		c.dbsClient = nil
		c.dbsClientErr = err
		return
	}
	c.dbsClient = client
	c.dbsClientErr = nil
}

func (c *Config) SetupProvisioningClient() {
	if c.iamClientErr != nil {
		c.provisioningClientErr = fmt.Errorf("IAM client error in SetupProvisioningClient: %w", c.iamClientErr)
		return
	}
	client, err := provisioning.NewClient(c.iamClient, &provisioning.Config{
		Region:      c.Region,
		Environment: c.Environment,
		DebugLog:    c.DebugWriter,
	})
	if err != nil {
		c.provisioningClient = nil
		c.provisioningClientErr = err
		return
	}
	c.provisioningClient = client
	c.provisioningClientErr = nil
}
