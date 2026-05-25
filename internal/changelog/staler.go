package changelog

import (
	"time"
)

// Staler flags entries whose published date is older than a configured
// staleness threshold. Flagged entries receive a "stale" tag so that
// downstream stages (notifiers, badgers, etc.) can handle or suppress them.
type Staler struct {
	threshold time.Duration
}

// NewStaler returns a Staler that marks entries older than threshold as stale.
// A zero or negative threshold disables staleness detection (no entries are
// ever flagged).
func NewStaler(threshold time.Duration) *Staler {
	return &Staler{threshold: threshold}
}

// Apply iterates over entries and appends a "stale" tag to any entry whose
// Date is non-zero and older than now minus the configured threshold.
// Entries with a zero Date are left untouched.
func (s *Staler) Apply(entries []Entry) []Entry {
	if s.threshold <= 0 {
		return entries
	}

	cutoff := time.Now().UTC().Add(-s.threshold)
	out := make([]Entry, len(entries))

	for i, e := range entries {
		if !e.Date.IsZero() && e.Date.Before(cutoff) {
			e.Tags = appendUniqueTag(e.Tags, "stale")
		}
		out[i] = e
	}

	return out
}

// appendUniqueTag appends tag to tags only if it is not already present.
func appendUniqueTag(tags []string, tag string) []string {
	for _, t := range tags {
		if t == tag {
			return tags
		}
	}
	return append(tags, tag)
}
