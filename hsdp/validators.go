package hsdp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	creds "github.com/philips-software/go-hsdp-api/s3creds"
	"github.com/robfig/cron/v3"
)

func validateCron(val interface{}, _ cty.Path) diag.Diagnostics {
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

var thresholdMapping = map[string]string{
	"cpu":          "threshold_cpu",
	"memory":       "threshold_memory",
	"http-rate":    "threshold_http_rate",
	"http-latency": "threshold_http_latency",
}
