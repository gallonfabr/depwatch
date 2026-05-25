package changelog

import "time"

// PinnedEntry records a dependency version that has been explicitly pinned
// by the operator. Pinned entries are excluded from digest notifications
// until the pin is lifted or the version advances past the pinned value.
type PinnedEntry struct {
	Dependency string
	Version    string
	PinnedAt   time.Time
	Reason     string
}

// Pinner filters out changelog entries whose (Dependency, Version) pair
// matches a pinned entry.
type Pinner struct {
	pins map[string]string // dependency -> pinned version
}

// NewPinner creates a Pinner pre-loaded with the supplied pinned entries.
// Passing a nil or empty slice produces a no-op Pinner.
func NewPinner(pins []PinnedEntry) *Pinner {
	m := make(map[string]string, len(pins))
	for _, p := range pins {
		if p.Dependency != "" && p.Version != "" {
			m[p.Dependency] = p.Version
		}
	}
	return &Pinner{pins: m}
}

// Apply removes any Entry whose Dependency is pinned at exactly the same
// Version string. Entries for un-pinned dependencies pass through unchanged.
func (p *Pinner) Apply(entries []Entry) []Entry {
	if len(p.pins) == 0 {
		return entries
	}
	out := entries[:0:0]
	for _, e := range entries {
		if pinnedVer, ok := p.pins[e.Dependency]; ok && pinnedVer == e.Version {
			continue
		}
		out = append(out, e)
	}
	return out
}

// IsPinned reports whether the given dependency+version pair is currently
// pinned.
func (p *Pinner) IsPinned(dependency, version string) bool {
	v, ok := p.pins[dependency]
	return ok && v == version
}
