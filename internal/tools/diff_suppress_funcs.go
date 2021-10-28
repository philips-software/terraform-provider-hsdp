package tools

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SuppressMulti(fns ...schema.SchemaDiffSuppressFunc) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		for _, f := range fns {
			if f(k, old, new, d) {
				return true
			}
		}
		return false
	}
}

func SuppressCaseDiffs(k, old, new string, d *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

func SuppressDefault(k, old, new string, d *schema.ResourceData) bool {
	if old == "default" && new == "" {
		return true
	}
	return false
}

func SuppressDefaultCommunicationChannel(k, old, new string, d *schema.ResourceData) bool {
	if (old == "email" || old == "sms") && new == "" {
		return true
	}
	return false
}

func SuppressEmptyPreferredLanguage(k, old, new string, d *schema.ResourceData) bool {
	if old != "" && new == "" {
		return true
	}
	return false
}

func SuppressWhenGenerated(k, old, new string, d *schema.ResourceData) bool {
	return new == ""
}

func SuppressEqualTimeOrMissing(k, old, new string, d *schema.ResourceData) bool {
	if new == "" { // Not set by us
		return true
	}
	oldTime, err := time.Parse(time.RFC3339, old)
	if err != nil {
		return false
	}
	newTime, err := time.Parse("2006-01-02", new)
	if err != nil {
		return false
	}
	return oldTime.Equal(newTime)
}
