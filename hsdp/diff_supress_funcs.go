package hsdp

import (
	"encoding/json"

	"github.com/hashicorp/terraform/helper/schema"
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
