package acc

import (
	"os"
	"sync"
	"testing"

	"github.com/philips-software/terraform-provider-hsdp/hsdp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// ProviderName for single configuration testing
	ProviderName = "hsdp"

	ResourcePrefix = "tf-acc-test"
)

const RFC3339RegexPattern = `^[0-9]{4}-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9](\.[0-9]+)?([Zz]|([+-]([01][0-9]|2[0-3]):[0-5][0-9]))$`

// Skip implements a wrapper for (*testing.T).Skip() to prevent unused linting reports
//
// Reference: https://github.com/dominikh/go-tools/issues/633#issuecomment-606560616
func Skip(t *testing.T, message string) {
	t.Skip(message)
}

// ProviderFactories is a static map containing only the main provider instance
//
// Use other ProviderFactories functions, such as FactoriesAlternate,
// for tests requiring special provider configurations.
var ProviderFactories map[string]func() (*schema.Provider, error)

// Provider is the "main" provider instance
//
// This Provider can be used in testing code for API calls without requiring
// the use of saving and referencing specific ProviderFactories instances.
//
// PreCheck(t) must be called before using this provider instance.
var Provider *schema.Provider

// testAccProviderConfigure ensures Provider is only configured once
//
// The PreCheck(t) function is invoked for every test and this prevents
// extraneous reconfiguration to the same values each time. However, this does
// not prevent reconfiguration that may happen should the address of
// Provider be erroneously reused in ProviderFactories.
var testAccProviderConfigure sync.Once

func init() {
	Provider = hsdp.Provider("test")

	// Always allocate a new provider instance each invocation, otherwise gRPC
	// ProviderConfigure() can overwrite configuration during concurrent testing.
	ProviderFactories = map[string]func() (*schema.Provider, error){
		ProviderName: func() (*schema.Provider, error) { return hsdp.Provider("test"), nil }, //nolint:unparam
	}
}

// PreCheck verifies and sets required provider testing configuration
//
// This PreCheck function should be present in every acceptance test. It allows
// test configurations to omit a provider configuration with region and ensures
// testing functions that attempt to call AWS APIs are previously configured.
//
// These verifications and configuration are preferred at this level to prevent
// provider developers from experiencing less clear errors for every test.
func PreCheck(t *testing.T) {
	// Since we are outside the scope of the Terraform configuration we must
	// call Configure() to properly initialize the provider configuration.
	testAccProviderConfigure.Do(func() {
		// TODO: add additional pre-checks here
		if AccUserGUID() == "" {
			t.Fatalf("HSDP_IAM_ACC_USER_GUID must be set")
		}
		if AccIAMOrgGUID() == "" {
			t.Fatalf("HSDP_IAM_ACC_ORG_GUID must be set")
		}
		if AccCDRURL() == "" {
			t.Fatalf("HSDP_CDR_ACC_URL must be set")
		}
	})
}

func AccUserGUID() string {
	return os.Getenv("HSDP_IAM_ACC_USER_GUID")
}

func AccIAMOrgGUID() string {
	return os.Getenv("HSDP_IAM_ACC_ORG_GUID")
}

func AccCDRURL() string {
	return os.Getenv("HSDP_CDR_ACC_URL")
}

func AccMDMClientID() string {
	return os.Getenv("HSDP_MDM_ACC_CLIENT_ID")
}

func AccMDMClientSecret() string {
	return os.Getenv("HSDP_MDM_ACC_CLIENT_SECRET")
}

func AccMDMOrgID() string {
	return os.Getenv("HSDP_MDM_ACC_ORG_ID")
}
