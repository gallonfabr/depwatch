package changelog

import (
	"testing"
	"time"
)

func TestNewThrottle_DefaultDelay(t *testing.T) {
	th := NewThrottle(0)
	if th.minDelay != time.Minute {
		t.Fatalf("expected default delay 1m, got %v", th.minDelay)
	}
}

func TestNewThrottle_CustomDelay(t *testing.T) {
	th := NewThrottle(5 * time.Second)
	if th.minDelay != 5*time.Second {
		t.Fatalf("expected 5s delay, got %v", th.minDelay)
	}
}

func TestThrottle_Allow_FirstCallPermitted(t *testing.T) {
	th := NewThrottle(10 * time.Second)
	if !th.Allow("dep-a") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestThrottle_Allow_SecondCallBlocked(t *testing.T) {
	th := NewThrottle(10 * time.Second)
	th.Allow("dep-a")
	if th.Allow("dep-a") {
		t.Fatal("expected second immediate call to be blocked")
	}
}

func TestThrottle_Allow_IndependentKeys(t *testing.T) {
	th := NewThrottle(10 * time.Second)
	th.Allow("dep-a")
	if !th.Allow("dep-b") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestThrottle_Allow_AfterDelayPermitted(t *testing.T) {
	th := NewThrottle(10 * time.Millisecond)
	th.Allow("dep-a")
	time.Sleep(20 * time.Millisecond)
	if !th.Allow("dep-a") {
		t.Fatal("expected call after delay to be allowed")
	}
}

func TestThrottle_Reset_AllowsImmediateFetch(t *testing.T) {
	th := NewThrottle(10 * time.Second)
	th.Allow("dep-a")
	th.Reset("dep-a")
	if !th.Allow("dep-a") {
		t.Fatal("expected allow after reset")
	}
}

func TestThrottle_Len_TracksKeys(t *testing.T) {
	th := NewThrottle(10 * time.Second)
	if th.Len() != 0 {
		t.Fatal("expected empty throttle")
	}
	th.Allow("dep-a")
	th.Allow("dep-b")
	if th.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", th.Len())
	}
}

func TestThrottle_Reset_UnknownKeyNoOp(t *testing.T) {
	th := NewThrottle(10 * time.Second)
	// Should not panic
	th.Reset("nonexistent")
	if th.Len() != 0 {
		t.Fatal("expected len 0 after reset of unknown key")
	}
}
