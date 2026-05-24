package changelog

import (
	"context"
	"errors"
	"time"
)

// RetryConfig holds configuration for the retry policy.
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// DefaultRetryConfig returns a RetryConfig with sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   500 * time.Millisecond,
		MaxDelay:    10 * time.Second,
	}
}

// Retry executes fn up to cfg.MaxAttempts times, backing off exponentially
// between attempts. It stops early if ctx is cancelled or fn returns a
// non-retryable error wrapped with ErrPermanent.
func Retry(ctx context.Context, cfg RetryConfig, fn func() error) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}

	delay := cfg.BaseDelay
	var lastErr error

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		var perm *PermanentError
		if errors.As(lastErr, &perm) {
			return lastErr
		}

		if attempt < cfg.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay *= 2
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		}
	}

	return lastErr
}

// PermanentError wraps an error to signal that retrying is futile.
type PermanentError struct {
	Cause error
}

func (e *PermanentError) Error() string {
	return "permanent: " + e.Cause.Error()
}

func (e *PermanentError) Unwrap() error { return e.Cause }

// Permanent wraps err so that Retry will not attempt further retries.
func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return &PermanentError{Cause: err}
}
