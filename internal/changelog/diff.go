package changelog

import "time"

// Diff compares two slices of Entry and returns entries that are new in
// current (i.e. present in current but absent from previous), matched by
// dependency name and version string.
type Diff struct{}

// NewDiff creates a new Diff instance.
func NewDiff() *Diff {
	return &Diff{}
}

// Apply returns entries from current that do not appear in previous.
// Comparison is performed on (Dependency, Version) pairs.
func (d *Diff) Apply(previous, current []Entry) []Entry {
	seen := make(map[string]struct{}, len(previous))
	for _, e := range previous {
		seen[diffKey(e)] = struct{}{}
	}

	var novel []Entry
	for _, e := range current {
		if _, exists := seen[diffKey(e)]; !exists {
			novel = append(novel, e)
		}
	}
	return novel
}

// Summary holds high-level statistics produced by Summarise.
type DiffSummary struct {
	Added   int
	Oldest  time.Time
	Newest  time.Time
}

// Summarise returns a DiffSummary for a slice of novel entries.
func (d *Diff) Summarise(entries []Entry) DiffSummary {
	if len(entries) == 0 {
		return DiffSummary{}
	}
	s := DiffSummary{Added: len(entries), Oldest: entries[0].Date, Newest: entries[0].Date}
	for _, e := range entries[1:] {
		if !e.Date.IsZero() && e.Date.Before(s.Oldest) {
			s.Oldest = e.Date
		}
		if e.Date.After(s.Newest) {
			s.Newest = e.Date
		}
	}
	return s
}

func diffKey(e Entry) string {
	return e.Dependency + "@" + e.Version
}
