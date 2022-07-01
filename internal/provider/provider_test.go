package provider

import (
	"testing"
)

func TestProvider(t *testing.T) {
	if err := Provider("v0.0.0").InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider("v0.0.0")
}
