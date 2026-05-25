package changelog

import (
	"errors"
	"sync"
	"time"
)

// ErrQuotaExceeded is returned when a dependency has exceeded its entry quota
// for the current period.
var ErrQuotaExceeded = errors.New("changelog: quota exceeded for dependency")

// Quota enforces a maximum number of entries per dependency within a rolling
// time window. Entries beyond the cap are silently dropped.
type Quota struct {
	mu      sync.Mutex
	counts  map[string][]time.Time
	max     int
	window  time.Duration
}

// NewQuota creates a Quota that allows at most max entries per dependency
// within the given window duration. Panics if max < 1 or window <= 0.
func NewQuota(max int, window time.Duration) *Quota {
	if max < 1 {
		panic("changelog: quota max must be at least 1")
	}
	if window <= 0 {
		panic("changelog: quota window must be positive")
	}
	return &Quota{
		counts: make(map[string][]time.Time),
		max:    max,
		window: window,
	}
}

// Apply filters entries, dropping any that exceed the per-dependency quota
// within the rolling window.
func (q *Quota) Apply(entries []Entry) []Entry {
	now := time.Now()
	q.mu.Lock()
	defer q.mu.Unlock()

	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		key := e.Dependency
		q.evict(key, now)
		if len(q.counts[key]) < q.max {
			q.counts[key] = append(q.counts[key], now)
			out = append(out, e)
		}
	}
	return out
}

// Reset clears all recorded counts, starting fresh for all dependencies.
func (q *Quota) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.counts = make(map[string][]time.Time)
}

// evict removes timestamps that have fallen outside the rolling window.
func (q *Quota) evict(key string, now time.Time) {
	times := q.counts[key]
	cutoff := now.Add(-q.window)
	i := 0
	for i < len(times) && times[i].Before(cutoff) {
		i++
	}
	q.counts[key] = times[i:]
}
