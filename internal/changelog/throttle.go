package changelog

import (
	"sync"
	"time"
)

// Throttle limits how frequently a source can be fetched by enforcing a
// minimum delay between successive calls for the same dependency key.
type Throttle struct {
	mu       sync.Mutex
	lastSeen map[string]time.Time
	minDelay time.Duration
}

// NewThrottle creates a Throttle with the given minimum delay between fetches.
// If minDelay is zero or negative it defaults to 1 minute.
func NewThrottle(minDelay time.Duration) *Throttle {
	if minDelay <= 0 {
		minDelay = time.Minute
	}
	return &Throttle{
		lastSeen: make(map[string]time.Time),
		minDelay: minDelay,
	}
}

// Allow reports whether a fetch for the given key is permitted at now.
// If allowed it records now as the last fetch time for that key.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	last, ok := t.lastSeen[key]
	if ok && time.Since(last) < t.minDelay {
		return false
	}
	t.lastSeen[key] = time.Now()
	return true
}

// Reset clears the recorded fetch time for the given key, allowing an
// immediate fetch on the next call to Allow.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastSeen, key)
}

// Len returns the number of keys currently tracked by the throttle.
func (t *Throttle) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.lastSeen)
}
