package tools

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/robfig/cron/v3"
)

func ValidateCron(val interface{}, _ cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	v, ok := val.(string)
	if !ok {
		return diag.FromErr(fmt.Errorf("string expected for CRON entry"))
	}
	_, err := cron.ParseStandard(v)
	if err != nil {
		return diag.FromErr(fmt.Errorf("invalid CRON entry format: %w", err))
	}
	return diags
}

func ValidateUpperString(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	u := strings.ToUpper(v)
	if v != u {
		errs = append(errs, fmt.Errorf("%q must be in uppercase: %s -> %s", key, v, u))
	}
	return
}
