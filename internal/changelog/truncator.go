package changelog

// Truncator trims a slice of Entry values to a maximum count per dependency.
// It is useful as a final safeguard before entries are handed to the digest
// builder, ensuring that no single dependency floods the output.
type Truncator struct {
	maxPerDep int
}

// NewTruncator returns a Truncator that keeps at most maxPerDep entries for
// each unique dependency name. If maxPerDep is zero or negative the Truncator
// is a no-op and returns all entries unchanged.
func NewTruncator(maxPerDep int) *Truncator {
	return &Truncator{maxPerDep: maxPerDep}
}

// Apply implements Transformer. It preserves the relative order of entries
// while enforcing the per-dependency cap.
func (t *Truncator) Apply(entries []Entry) []Entry {
	if t.maxPerDep <= 0 {
		return entries
	}

	counts := make(map[string]int)
	out := make([]Entry, 0, len(entries))

	for _, e := range entries {
		if counts[e.Dependency] < t.maxPerDep {
			out = append(out, e)
			counts[e.Dependency]++
		}
	}

	return out
}
