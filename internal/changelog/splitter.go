package changelog

// Splitter partitions a flat slice of entries into per-dependency buckets.
// Each bucket preserves the original relative order of its entries.
type Splitter struct{}

// NewSplitter returns a new Splitter.
func NewSplitter() *Splitter {
	return &Splitter{}
}

// Split groups entries by their Dependency field and returns a map whose
// keys are dependency names and whose values are the matching entries.
// Entries with an empty Dependency field are placed under the empty-string key.
func (s *Splitter) Split(entries []Entry) map[string][]Entry {
	result := make(map[string][]Entry)
	for _, e := range entries {
		result[e.Dependency] = append(result[e.Dependency], e)
	}
	return result
}

// Keys returns the dependency names present in entries, in first-seen order.
func (s *Splitter) Keys(entries []Entry) []string {
	seen := make(map[string]struct{})
	var keys []string
	for _, e := range entries {
		if _, ok := seen[e.Dependency]; !ok {
			seen[e.Dependency] = struct{}{}
			keys = append(keys, e.Dependency)
		}
	}
	return keys
}
