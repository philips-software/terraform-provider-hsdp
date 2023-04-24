package provider_test

import (
	"testing"

	"github.com/philips-software/terraform-provider-hsdp/internal/provider"
)

func TestProvider(t *testing.T) {
	if err := provider.Provider("v0.0.0").InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = provider.Provider("v0.0.0")
}
