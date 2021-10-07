package hsdp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/philips-software/go-hsdp-api/iam"
)

func tryIAMCall(operation func() (*iam.Response, error), retryOnCodes ...int) error {
	if len(retryOnCodes) == 0 {
		retryOnCodes = []int{http.StatusUnprocessableEntity, http.StatusInternalServerError}
	}
	doOp := func() error {
		resp, err := operation()
		if err == nil {
			return nil
		}
		if resp == nil {
			return backoff.Permanent(fmt.Errorf("response was nil: %w", err))
		}
		shouldRetry := false
		for _, c := range retryOnCodes {
			if c == resp.StatusCode {
				shouldRetry = true
				break
			}
		}
		if shouldRetry {
			return err
		}
		return backoff.Permanent(err)
	}
	return backoff.Retry(doOp, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
}

// difference returns the elements in a that aren't in b
func difference(a, b []string) []string {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := []string{}
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab
}

func nextQuarterStart(now time.Time) time.Time {
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

func slidingExpiresOn(now time.Time) string {
	expiresOn := nextQuarterStart(now)
	delta := expiresOn.Sub(now).Hours() / 24
	if delta < 30 {
		expiresOn = nextQuarterStart(expiresOn)
	}
	return expiresOn.Format(time.RFC3339)
}
