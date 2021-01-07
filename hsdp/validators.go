package hsdp

import (
	"encoding/json"
	"fmt"
	"github.com/philips-software/go-hsdp-api/cartel"
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

func validateSubnet(client *cartel.Client, v string) error {
	subnets, _, err := client.GetAllSubnets()

	if err != nil {
		return err
	}
	var availableSubnets []string
	for _, subnet := range *subnets {
		availableSubnets = append(availableSubnets, subnet.ID)
		if v == subnet.ID {
			return nil
		}
	}
	return fmt.Errorf("unsupported subnet: %s (available: %s)", v, strings.Join(availableSubnets, ", "))
}

var thresholdMapping = map[string]string{
	"cpu":          "threshold_cpu",
	"memory":       "threshold_memory",
	"http-rate":    "threshold_http_rate",
	"http-latency": "threshold_http_latency",
}
