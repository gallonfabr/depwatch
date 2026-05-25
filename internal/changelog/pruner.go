package changelog

import "time"

// Pruner removes entries that fall outside a retention window defined by a
// maximum age. Entries with a zero date are always kept.
type Pruner struct {
	maxAge time.Duration
}

// PrunerOption configures a Pruner.
type PrunerOption func(*Pruner)

// WithMaxAge sets the maximum age for retained entries.
func WithMaxAge(d time.Duration) PrunerOption {
	return func(p *Pruner) {
		if d > 0 {
			p.maxAge = d
		}
	}
}

// NewPruner returns a Pruner with the supplied options.
// The default max age is 30 days.
func NewPruner(opts ...PrunerOption) *Pruner {
	p := &Pruner{maxAge: 30 * 24 * time.Hour}
	for _, o := range opts {
		o(p)
	}
	return p
}

// Apply removes entries older than the configured max age.
// Entries with a zero date are passed through unchanged.
func (p *Pruner) Apply(entries []Entry) []Entry {
	if len(entries) == 0 {
		return entries
	}
	cutoff := time.Now().UTC().Add(-p.maxAge)
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if e.Date.IsZero() || !e.Date.Before(cutoff) {
			out = append(out, e)
		}
	}
	return out
}
