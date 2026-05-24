package changelog

import (
	"context"
	"testing"
	"time"
)

func TestNewRateLimiter_InvalidMaxCalls(t *testing.T) {
	_, err := NewRateLimiter(0, time.Second)
	if err == nil {
		t.Fatal("expected error for maxCalls=0")
	}
}

func TestNewRateLimiter_InvalidWindow(t *testing.T) {
	_, err := NewRateLimiter(1, 0)
	if err == nil {
		t.Fatal("expected error for window=0")
	}
}

func TestNewRateLimiter_Valid(t *testing.T) {
	rl, err := NewRateLimiter(5, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rl == nil {
		t.Fatal("expected non-nil RateLimiter")
	}
}

func TestRateLimiter_Allow_UnderLimit(t *testing.T) {
	rl, _ := NewRateLimiter(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !rl.Allow("dep-a") {
			t.Fatalf("call %d should be allowed", i+1)
		}
	}
}

func TestRateLimiter_Allow_ExceedsLimit(t *testing.T) {
	rl, _ := NewRateLimiter(2, time.Minute)
	rl.Allow("dep-b")
	rl.Allow("dep-b")
	if rl.Allow("dep-b") {
		t.Fatal("third call should be blocked")
	}
}

func TestRateLimiter_Allow_IndependentKeys(t *testing.T) {
	rl, _ := NewRateLimiter(1, time.Minute)
	if !rl.Allow("dep-x") {
		t.Fatal("first call for dep-x should be allowed")
	}
	if !rl.Allow("dep-y") {
		t.Fatal("first call for dep-y should be allowed")
	}
	if rl.Allow("dep-x") {
		t.Fatal("second call for dep-x should be blocked")
	}
}

func TestRateLimiter_Reset_ClearsKey(t *testing.T) {
	rl, _ := NewRateLimiter(1, time.Minute)
	rl.Allow("dep-c")
	rl.Reset("dep-c")
	if !rl.Allow("dep-c") {
		t.Fatal("call after reset should be allowed")
	}
}

func TestRateLimiter_Allow_ExpiredWindowCleared(t *testing.T) {
	rl, _ := NewRateLimiter(1, 50*time.Millisecond)
	rl.Allow("dep-d")
	time.Sleep(60 * time.Millisecond)
	if !rl.Allow("dep-d") {
		t.Fatal("call after window expiry should be allowed")
	}
}

func TestRateLimiter_Wait_ContextCancelled(t *testing.T) {
	rl, _ := NewRateLimiter(1, time.Second)
	rl.Allow("dep-e") // exhaust limit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err := rl.Wait(ctx, "dep-e")
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}
