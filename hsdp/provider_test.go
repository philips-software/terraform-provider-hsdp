package hsdp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider("v0.0.0")
	testAccProviders = map[string]*schema.Provider{
		"hsdp": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider("v0.0.0").InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider("v0.0.0")
}

func testAccPreCheck(t *testing.T) {
}
