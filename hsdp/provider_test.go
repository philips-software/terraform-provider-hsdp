package hsdp_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// credentialAttributeRe matches string attribute names that look like they hold
// credential material and are expected to carry Sensitive: true.
var credentialAttributeRe = regexp.MustCompile(`(?i)(secret|password|token|_key)$`)

// credentialAttributeAllowList covers identifier-style names that match
// credentialAttributeRe but are not secrets themselves (e.g. an ID paired with a key).
var credentialAttributeAllowList = map[string]bool{
	"client_id":                    true,
	"hsdp_product_key":             true,
	"key_id":                       true,
	"public_key":                   true,
	"kinesis_stream_partition_key": true,
}

func TestProvider_credentialAttributesAreSensitive(t *testing.T) {
	p := hsdp.Provider("v0.0.0")

	var walk func(prefix string, res *schema.Resource)
	walk = func(prefix string, res *schema.Resource) {
		if res == nil {
			return
		}
		for name, s := range res.Schema {
			path := prefix + "." + name
			if s.Type == schema.TypeString && credentialAttributeRe.MatchString(name) &&
				!credentialAttributeAllowList[name] && !s.Sensitive {
				t.Errorf("%s: attribute %q looks like a credential but is missing Sensitive: true", path, name)
			}
			if s.Elem != nil {
				if elemRes, ok := s.Elem.(*schema.Resource); ok {
					walk(path, elemRes)
				}
			}
		}
	}

	for name, res := range p.ResourcesMap {
		walk("resource."+name, res)
	}
	for name, res := range p.DataSourcesMap {
		walk("data."+name, res)
	}
}
