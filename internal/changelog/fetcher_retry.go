package changelog

import (
	"context"
	"fmt"
	"time"
)

// RetryFetcher wraps any Fetcher and retries transient failures according to
// a RetryConfig.
type RetryFetcher struct {
	inner  Fetcher
	cfg    RetryConfig
}

// NewRetryFetcher returns a RetryFetcher that delegates to inner and retries
// up to cfg.MaxAttempts times on transient errors.
func NewRetryFetcher(inner Fetcher, cfg RetryConfig) (*RetryFetcher, error) {
	if inner == nil {
		return nil, fmt.Errorf("changelog: RetryFetcher requires a non-nil inner Fetcher")
	}
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}
	if cfg.BaseDelay <= 0 {
		cfg.BaseDelay = 500 * time.Millisecond
	}
	if cfg.MaxDelay <= 0 {
		cfg.MaxDelay = 10 * time.Second
	}
	return &RetryFetcher{inner: inner, cfg: cfg}, nil
}

// Fetch calls the inner Fetcher's Fetch method, retrying on transient errors.
func (r *RetryFetcher) Fetch(url string) (string, error) {
	var result string
	err := Retry(context.Background(), r.cfg, func() error {
		var fetchErr error
		result, fetchErr = r.inner.Fetch(url)
		return fetchErr
	})
	return result, err
}
