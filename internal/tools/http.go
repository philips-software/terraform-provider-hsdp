package tools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cenkalti/backoff/v4"
)

var (
	StandardRetryOnCodes = []int{http.StatusForbidden, http.StatusInternalServerError, http.StatusServiceUnavailable, http.StatusTooManyRequests, http.StatusBadGateway, http.StatusGatewayTimeout}
)

func TryHTTPCall(ctx context.Context, numberOfTries uint64, operation func() (*http.Response, error), retryOnCodes ...int) error {
	if len(retryOnCodes) == 0 {
		retryOnCodes = StandardRetryOnCodes
	}
	count := 0
	doOp := func() error {
		resp, err := operation()
		if err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return backoff.Permanent(fmt.Errorf("context was cancelled: %w", err))
		default:
		}
		shouldRetry := false
		if resp == nil {
			err = fmt.Errorf("response was nil: %w", err)
			shouldRetry = true
		}
		if resp != nil {
			for _, c := range retryOnCodes {
				if c == resp.StatusCode {
					shouldRetry = true
					break
				}
			}
		}
		if shouldRetry {
			count = count + 1
			httpCode := 0
			if resp != nil {
				httpCode = resp.StatusCode
			}
			return fmt.Errorf("retry %d due to HTTP %d: %w", count, httpCode, err)
		}
		return backoff.Permanent(fmt.Errorf("retry %d permanent: %w", count, err))
	}
	return backoff.Retry(doOp, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), numberOfTries))
}
