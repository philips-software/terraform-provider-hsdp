package tools

import (
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func ReverseString(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}

// Difference returns the elements in a that aren't in b
func Difference(a, b []string) []string {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	var ab []string
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab
}

func NextQuarterStart(now time.Time) time.Time {
	year := now.Year()
	month := now.Month()
	if now.Day() > 1 {
		month += 1
	}
	month += 4 - (month % 3)
	if month > 12 {
		year += 1
	}
	month = month % 12
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

func SlidingExpiresOn(now time.Time) string {
	expiresOn := NextQuarterStart(now)
	delta := expiresOn.Sub(now).Hours() / 24
	if delta < 30 {
		expiresOn = NextQuarterStart(expiresOn)
	}
	return expiresOn.Format(time.RFC3339)
}

func SSHAgentReachable() bool {
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func PrunePorts(i []int, pruneList []int) []int {
	// Sort
	ports := i
	sort.Ints(ports)
	// Prune
	j := 0
	for _, v := range ports {
		prune := false
		for _, p := range pruneList {
			if v == p {
				prune = true
				continue
			}
		}
		if prune {
			continue
		}
		ports[j] = v
		j++
	}

	return ports[:j]
}

// ExpandStringList takes the result of flatmap.Expand for an array of strings
// and returns a []string
func ExpandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

func ContainsString(haystack []string, needle string) bool {
	for _, a := range haystack {
		if strings.EqualFold(a, needle) {
			return true
		}
	}
	return false
}

func ExpandIntList(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(int)
		if ok && val != 0 {
			vs = append(vs, val)
		}
	}
	return vs
}

func CollectList(fieldName string, d *schema.ResourceData) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	list := d.Get(fieldName).([]interface{})
	commands := make([]string, 0)
	for i := 0; i < len(list); i++ {
		commands = append(commands, list[i].(string))
	}
	return commands, diags

}

func CheckForIAMPermissionErrors(client iam.TokenRefresher, resp *http.Response, err error) error {
	if resp == nil || resp.StatusCode > 500 {
		return err
	}
	if resp.StatusCode == http.StatusForbidden {
		_ = client.TokenRefresh()
		return err
	}
	return backoff.Permanent(err)
}

// DisableFHIRValidation disable validation for the request
func DisableFHIRValidation(request *http.Request) error {
	request.Header.Set("X-Validate-Resource", "false")
	return nil
}

func String(str string) *string {
	return &str
}
