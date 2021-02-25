package hsdp

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	creds "github.com/philips-software/go-hsdp-api/s3creds"
)

func suppressEquivalentPolicyDiffs(k, old, new string, d *schema.ResourceData) bool {
	oldPolicy := &creds.Policy{}
	if err := json.Unmarshal([]byte(old), oldPolicy); err != nil {
		return false
	}
	newPolicy := &creds.Policy{}
	if err := json.Unmarshal([]byte(new), newPolicy); err != nil {
		return false
	}
	return oldPolicy.Equals(newPolicy)
}

func suppressCaseDiffs(k, old, new string, d *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

func suppressDefault(k, old, new string, d *schema.ResourceData) bool {
	if old == "default" && new == "" {
		return true
	}
	return false
}

func suppressWhenGenerated(k, old, new string, d *schema.ResourceData) bool {
	if new == "" {
		return true
	}
	return false
}
