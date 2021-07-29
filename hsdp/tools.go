package hsdp

import (
	"fmt"
	"net/http"

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
