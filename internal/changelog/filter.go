package changelog

import "time"

// Entry represents a single parsed changelog entry.
type Entry struct {
	Version string
	Date    time.Time
	Body    string
}

// Filter holds criteria for filtering changelog entries.
type Filter struct {
	Since time.Time
	Limit int
}

// NewFilter creates a Filter that returns entries newer than since.
// If limit is <= 0, no limit is applied.
func NewFilter(since time.Time, limit int) *Filter {
	return &Filter{Since: since, Limit: limit}
}

// Apply returns only the entries that match the filter criteria.
// Entries are assumed to be ordered newest-first.
func (f *Filter) Apply(entries []Entry) []Entry {
	var result []Entry
	for _, e := range entries {
		if !f.Since.IsZero() && !e.Date.After(f.Since) {
			continue
		}
		result = append(result, e)
		if f.Limit > 0 && len(result) >= f.Limit {
			break
		}
	}
	return result
}

// FilterNew returns entries whose version is not present in seen.
// seen is a set of version strings already delivered to the user.
func FilterNew(entries []Entry, seen map[string]bool) []Entry {
	var result []Entry
	for _, e := range entries {
		if !seen[e.Version] {
			result = append(result, e)
		}
	}
	return result
}
