package changelog

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RateLimiter enforces a maximum number of calls per time window per key.
type RateLimiter struct {
	mu       sync.Mutex
	window   time.Duration
	maxCalls int
	buckets  map[string][]time.Time
}

// NewRateLimiter creates a RateLimiter that allows at most maxCalls within window.
// Returns an error if maxCalls < 1 or window <= 0.
func NewRateLimiter(maxCalls int, window time.Duration) (*RateLimiter, error) {
	if maxCalls < 1 {
		return nil, fmt.Errorf("maxCalls must be at least 1, got %d", maxCalls)
	}
	if window <= 0 {
		return nil, fmt.Errorf("window must be positive, got %s", window)
	}
	return &RateLimiter{
		window:   window,
		maxCalls: maxCalls,
		buckets:  make(map[string][]time.Time),
	}, nil
}

// Allow reports whether a call for the given key is permitted under the rate limit.
// It prunes expired timestamps and records the current call if allowed.
func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-r.window)

	times := r.buckets[key]
	valid := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= r.maxCalls {
		r.buckets[key] = valid
		return false
	}

	r.buckets[key] = append(valid, now)
	return true
}

// Wait blocks until a call for the given key is permitted or ctx is cancelled.
func (r *RateLimiter) Wait(ctx context.Context, key string) error {
	for {
		if r.Allow(key) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(r.window / time.Duration(r.maxCalls)):
		}
	}
}

// Reset clears all recorded call timestamps for the given key.
func (r *RateLimiter) Reset(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.buckets, key)
}
