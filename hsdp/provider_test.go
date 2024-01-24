package hsdp_test

import (
	"testing"

	"github.com/philips-software/terraform-provider-hsdp/hsdp"
)

func TestProvider(t *testing.T) {
	if err := hsdp.Provider("v0.0.0").InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = hsdp.Provider("v0.0.0")
}
