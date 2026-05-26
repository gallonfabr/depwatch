package changelog

// CrossSourceDeduplicator removes entries that appear in more than one
// upstream source feed, keeping only the first occurrence encountered
// across successive Apply calls.
//
// It is distinct from Deduplicator (which deduplicates within a single
// dependency's own history) and CrossDeduplicator (which works on
// version+body keys). CrossSourceDeduplicator keys on dependency+version
// so that the same release fetched from both an HTTP changelog and a
// GitHub releases feed is emitted only once.
type CrossSourceDeduplicator struct {
	seen map[string]struct{}
}

// NewCrossSourceDeduplicator returns an initialised CrossSourceDeduplicator.
func NewCrossSourceDeduplicator() *CrossSourceDeduplicator {
	return &CrossSourceDeduplicator{seen: make(map[string]struct{})}
}

func crossSourceKey(e Entry) string {
	return e.Dependency + "\x00" + e.Version
}

// Apply filters out any entry whose dependency+version pair has already
// been seen in a previous or the current Apply call.
func (c *CrossSourceDeduplicator) Apply(entries []Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		k := crossSourceKey(e)
		if _, exists := c.seen[k]; exists {
			continue
		}
		c.seen[k] = struct{}{}
		out = append(out, e)
	}
	return out
}

// Reset clears the internal seen-set so the deduplicator can be reused
// across digest cycles.
func (c *CrossSourceDeduplicator) Reset() {
	c.seen = make(map[string]struct{})
}
