package tools

import (
	"fmt"
	"strings"

	cfg "github.com/philips-software/go-dip-api/config"
)

func ValidateUpperString(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	u := strings.ToUpper(v)
	if v != u {
		errs = append(errs, fmt.Errorf("%q must be in uppercase: %s -> %s", key, v, u))
	}
	return
}

func ValidateRegion(i interface{}, k string) (warns []string, es []error) {
	r := i.(string)
	if r == "dev" { // dev is special case region so don't complain
		return
	}
	d, err := cfg.New(cfg.WithRegion(r))
	if err != nil {
		es = append(es, err)
	}
	if d.Service("cf").URL == "" {
		warns = append(warns, fmt.Sprintf("no Cloud foundry presence in region '%s'", r))
	}
	if d.Env("prod").Service("iam").URL == "" {
		warns = append(warns, fmt.Sprintf("no production IAM presence in region '%s'", r))
	}
	return
}

func ValidateEnvironment(i interface{}, k string) (warns []string, es []error) {
	env := i.(string)
	if !ContainsString([]string{"dev", "preview", "client-test", "prod", "production"}, env) {
		es = append(es, fmt.Errorf("environment '%s' is not a supported one", env))
	}
	return
}
