package changelog

// Booster promotes entries matching a set of dependency names to the top
// of the slice, preserving relative order within each group.
//
// Entries whose dependency name appears in the priority list are moved before
// all non-priority entries. Within each group the original order is kept.
type Booster struct {
	priority map[string]struct{}
}

// NewBooster returns a Booster that will promote entries belonging to any of
// the supplied dependency names. Duplicate names are silently ignored.
func NewBooster(deps ...string) *Booster {
	p := make(map[string]struct{}, len(deps))
	for _, d := range deps {
		if d != "" {
			p[d] = struct{}{}
		}
	}
	return &Booster{priority: p}
}

// Apply reorders entries so that priority dependencies appear first.
// The input slice is not mutated; a new slice is returned.
func (b *Booster) Apply(entries []Entry) []Entry {
	if len(entries) == 0 {
		return entries
	}

	prioritised := make([]Entry, 0, len(entries))
	rest := make([]Entry, 0, len(entries))

	for _, e := range entries {
		if _, ok := b.priority[e.Dependency]; ok {
			prioritised = append(prioritised, e)
		} else {
			rest = append(rest, e)
		}
	}

	return append(prioritised, rest...)
}
