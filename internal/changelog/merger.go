package changelog

// Merger combines multiple slices of Entry into a single deduplicated slice,
// preserving relative order from each source.
type Merger struct {
	dedup *Deduplicator
}

// NewMerger returns a Merger backed by a fresh Deduplicator.
func NewMerger() *Merger {
	return &Merger{dedup: NewDeduplicator()}
}

// Merge accepts a variadic number of Entry slices and returns a single slice
// that contains every entry exactly once, in the order they are first
// encountered (left-to-right, top-to-bottom).
func (m *Merger) Merge(sources ...[]Entry) []Entry {
	var combined []Entry
	for _, src := range sources {
		combined = append(combined, src...)
	}
	return m.dedup.Apply(combined)
}

// MergeAll is a package-level convenience wrapper around NewMerger().Merge.
// Each call creates a fresh Deduplicator, so cross-call state is not shared.
func MergeAll(sources ...[]Entry) []Entry {
	return NewMerger().Merge(sources...)
}
