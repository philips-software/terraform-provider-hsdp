package tenant_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/tenant"
	"github.com/stretchr/testify/assert"
)

// Unit test for the generateAPIKeyAndSignature function
func TestGenerateAPIKeyAndSignature(t *testing.T) {
	// Create a ResourceData instance for testing
	r := tenant.ResourceTenantKey()
	d := r.Data(nil)

	// Set test values
	_ = d.Set("project", "test-project")
	_ = d.Set("organization", "test-org")
	_ = d.Set("signing_key", "test-signing-key")
	_ = d.Set("scopes", []interface{}{"scope1", "scope2"})
	_ = d.Set("region", "us-east")
	_ = d.Set("environment", "prod")
	_ = d.Set("expiration", "2025-12-31T23:59:59Z")
	_ = d.Set("salt", "test-salt-value")

	// Call the exported test function
	apiKey, signature, err := tenant.TestGenerateAPIKeyAndSignature(d)

	// Assert results
	assert.Nil(t, err)
	assert.NotEmpty(t, apiKey)
	assert.NotEmpty(t, signature)
}

// Acceptance test for the tenant key resource
func TestAccResourceTenantKey_basic(t *testing.T) {
	resourceName := "hsdp_tenant_key.test"
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	testSigningKey := fmt.Sprintf("test-signing-key-%s", randomName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTenantKey(randomName, testSigningKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project", "test-project-"+randomName),
					resource.TestCheckResourceAttr(resourceName, "organization", "test-org-"+randomName),
					resource.TestCheckResourceAttr(resourceName, "signing_key", testSigningKey),
					resource.TestCheckResourceAttr(resourceName, "region", "us-east"),
					resource.TestCheckResourceAttr(resourceName, "environment", "prod"),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^[A-Za-z0-9_-]+$`)),
					resource.TestMatchResourceAttr(resourceName, "signature", regexp.MustCompile(`^[A-Za-z0-9_-]+$`)),
				),
			},
		},
	})
}

func testAccResourceTenantKey(name, signingKey string) string {
	return fmt.Sprintf(`
resource "hsdp_tenant_key" "test" {
  project      = "test-project-%s"
  organization = "test-org-%s"
  signing_key  = "%s"
  scopes       = ["scope1", "scope2"]
  expiration   = "2025-12-31T23:59:59Z"
  salt         = "test-salt-%s"
  region       = "us-east"
  environment  = "prod"
}
`, name, name, signingKey, name)
}

// Mock test to ensure API key generation and validation works as expected
func TestAccResourceTenantKey_mock(t *testing.T) {
	r := tenant.ResourceTenantKey()
	d := r.Data(nil)

	_ = d.Set("project", "test-project")
	_ = d.Set("organization", "test-org")
	_ = d.Set("signing_key", "test-signing-key")
	_ = d.Set("scopes", []interface{}{"scope1", "scope2"})
	_ = d.Set("region", "us-east")
	_ = d.Set("environment", "prod")
	_ = d.Set("expiration", "2025-12-31T23:59:59Z")
	_ = d.Set("salt", "test-salt-value")

	// Test create operation (can be expanded if we can mock the keys.GenerateAPIKey function)
	t.Run("create operation", func(t *testing.T) {
		// This is a placeholder for more comprehensive testing when mocking is set up
		resource := tenant.ResourceTenantKey()
		createFunc := resource.CreateContext
		diags := createFunc(nil, d, nil)

		// Without proper mocking of the keys package, we expect errors
		assert.Equal(t, 0, len(diags))
	})

	// Test read operation
	t.Run("read operation", func(t *testing.T) {
		// This is a placeholder for more comprehensive testing when mocking is set up
		resource := tenant.ResourceTenantKey()
		readFunc := resource.ReadContext
		diags := readFunc(nil, d, nil)

		// Without proper mocking of the keys package, we expect errors
		assert.Equal(t, 0, len(diags))
	})

	// Test delete operation
	t.Run("delete operation", func(t *testing.T) {
		// Set an ID so we can verify it gets cleared
		d.SetId("test-id")

		// Call the resource's Delete function
		resource := tenant.ResourceTenantKey()
		deleteFunc := resource.DeleteContext
		diags := deleteFunc(nil, d, nil)

		assert.Equal(t, 0, len(diags))
		assert.Empty(t, d.Id())
	})
}

// Test time validation on the expiration field
func TestTenantKeyTimeValidation(t *testing.T) {
	r := tenant.ResourceTenantKey()
	schema := r.Schema["expiration"]
	validateFunc := schema.ValidateFunc

	// Test valid time formats
	validTimes := []string{
		"2025-12-31T23:59:59Z",
		"2025-01-01T00:00:00Z",
		"2030-06-15T12:30:45Z",
	}

	for _, timeStr := range validTimes {
		warns, errs := validateFunc(timeStr, "expiration")
		assert.Equal(t, 0, len(warns), "No warnings expected for valid time: %s", timeStr)
		assert.Equal(t, 0, len(errs), "No errors expected for valid time: %s", timeStr)
	}

	// Test invalid time formats
	invalidTimes := []string{
		"",
		"invalid",
		"2025-12-31",           // Missing time component
		"2025-13-31T23:59:59Z", // Invalid month
		"2025-12-32T23:59:59Z", // Invalid day
		"2025-12-31 23:59:59",  // Missing 'T' and 'Z'
		"25-12-31T23:59:59Z",   // Incomplete year
	}

	for _, timeStr := range invalidTimes {
		_, errs := validateFunc(timeStr, "expiration")
		assert.True(t, len(errs) > 0, "Expected errors for invalid time: %s", timeStr)
	}
}
