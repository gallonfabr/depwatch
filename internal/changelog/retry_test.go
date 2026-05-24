package changelog

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errTransient = errors.New("transient error")

func TestRetry_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := Retry(context.Background(), DefaultRetryConfig(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetry_RetriesOnTransientError(t *testing.T) {
	calls := 0
	cfg := RetryConfig{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond}
	err := Retry(context.Background(), cfg, func() error {
		calls++
		if calls < 3 {
			return errTransient
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRetry_ExhaustsAttempts(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 2, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond}
	calls := 0
	err := Retry(context.Background(), cfg, func() error {
		calls++
		return errTransient
	})
	if !errors.Is(err, errTransient) {
		t.Fatalf("expected transient error, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestRetry_StopsOnPermanentError(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 5, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond}
	calls := 0
	err := Retry(context.Background(), cfg, func() error {
		calls++
		return Permanent(errTransient)
	})
	var perm *PermanentError
	if !errors.As(err, &perm) {
		t.Fatalf("expected PermanentError, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetry_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	err := Retry(ctx, DefaultRetryConfig(), func() error {
		calls++
		return nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls != 0 {
		t.Fatalf("expected 0 calls, got %d", calls)
	}
}

// TestRetry_ContextCancelledDuringRetry verifies that a context cancelled
// between attempts causes Retry to stop and return context.Canceled.
func TestRetry_ContextCancelledDuringRetry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	cfg := RetryConfig{MaxAttempts: 5, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond}
	err := Retry(ctx, cfg, func() error {
		calls++
		if calls == 2 {
			cancel()
		}
		return errTransient
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestPermanent_NilPassthrough(t *testing.T) {
	if Permanent(nil) != nil {
		t.Fatal("expected nil for Permanent(nil)")
	}
}

func TestDefaultRetryConfig_Values(t *testing.T) {
	cfg := DefaultRetryConfig()
	if cfg.MaxAttempts != 3 {
		t.Fatalf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.BaseDelay != 500*time.Millisecond {
		t.Fatalf("unexpected BaseDelay: %v", cfg.BaseDelay)
	}
}
