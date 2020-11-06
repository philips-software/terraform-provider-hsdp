package hsdp

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	creds "github.com/philips-software/go-hsdp-api/credentials"
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
	if strings.ToLower(old) == strings.ToLower(new) {
		return true
	}
	return false
}

func suppressOnID(k, old, new string, d *schema.ResourceData) bool {
	if d.Id() != "" {
		return true
	}
	return false
}
