package hsdp

import (
	"fmt"
	"strings"
)

func validateUpperString(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	u := strings.ToUpper(v)
	if v != u {
		errs = append(errs, fmt.Errorf("%q must be in uppercase: %s -> %s", key, v, u))
	}
	return
}
