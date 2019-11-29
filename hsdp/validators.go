package hsdp

import (
	"encoding/json"
	"fmt"
	"strings"

	creds "github.com/philips-software/go-hsdp-api/credentials"
)

func validateUpperString(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	u := strings.ToUpper(v)
	if v != u {
		errs = append(errs, fmt.Errorf("%q must be in uppercase: %s -> %s", key, v, u))
	}
	return
}

func validatePolicyJSON(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	var policy creds.Policy

	err := json.Unmarshal([]byte(v), &policy)
	if err != nil {
		errs = append(errs, fmt.Errorf("%q contains invalid JSON: %s, %v", key, v, err))
	}
	return
}
