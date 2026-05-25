package changelog

import "strings"

// Deduplicator is already taken; this is a cross-source merger that detects
// near-duplicate entries by comparing normalised version strings and
// dependency names, keeping only the first occurrence.

// CrossDeduplicator removes entries that share the same dependency and
// version across multiple sources, preferring entries that appear earlier
// in the slice.
type CrossDeduplicator struct {
	seen map[string]struct{}
}

// NewCrossDeduplicator returns a ready-to-use CrossDeduplicator.
func NewCrossDeduplicator() *CrossDeduplicator {
	return &CrossDeduplicator{seen: make(map[string]struct{})}
}

// Apply filters entries, removing any whose (dependency, version) pair has
// already been observed in a previous or current call.
func (c *CrossDeduplicator) Apply(entries []Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		k := crossKey(e.Dependency, e.Version)
		if _, exists := c.seen[k]; exists {
			continue
		}
		c.seen[k] = struct{}{}
		out = append(out, e)
	}
	return out
}

// Reset clears the internal seen-set so the deduplicator can be reused for a
// fresh run without allocating a new instance.
func (c *CrossDeduplicator) Reset() {
	c.seen = make(map[string]struct{})
}

func crossKey(dep, version string) string {
	return strings.ToLower(dep) + "@" + strings.ToLower(version)
}
