package changelog

// Deduplicator removes duplicate changelog entries based on dependency name
// and version, ensuring each version is only reported once per run.
type Deduplicator struct {
	seen map[string]struct{}
}

// NewDeduplicator creates a new Deduplicator with an empty seen set.
func NewDeduplicator() *Deduplicator {
	return &Deduplicator{
		seen: make(map[string]struct{}),
	}
}

// Apply filters out entries that have already been seen, returning only
// entries with unique dependency+version combinations. Seen entries are
// recorded so subsequent calls also exclude them.
func (d *Deduplicator) Apply(dep string, entries []Entry) []Entry {
	if len(entries) == 0 {
		return entries
	}

	unique := make([]Entry, 0, len(entries))
	for _, e := range entries {
		key := dep + "@" + e.Version
		if _, exists := d.seen[key]; exists {
			continue
		}
		d.seen[key] = struct{}{}
		unique = append(unique, e)
	}
	return unique
}

// Reset clears all recorded entries, allowing previously seen versions to
// be reported again.
func (d *Deduplicator) Reset() {
	d.seen = make(map[string]struct{})
}

// Len returns the number of unique dep@version keys currently tracked.
func (d *Deduplicator) Len() int {
	return len(d.seen)
}
