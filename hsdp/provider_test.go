package hsdp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider("v0.0.0").(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"hsdp": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider("v0.0.0").(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider("v0.0.0")
}

func testAccPreCheck(t *testing.T) {
}
