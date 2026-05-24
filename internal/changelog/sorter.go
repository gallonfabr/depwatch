package changelog

import (
	"sort"
	"time"
)

// SortOrder defines the ordering direction for changelog entries.
type SortOrder int

const (
	// SortDescending orders entries newest-first (default).
	SortDescending SortOrder = iota
	// SortAscending orders entries oldest-first.
	SortAscending
)

// Sorter sorts changelog entries by date.
type Sorter struct {
	order SortOrder
}

// NewSorter returns a Sorter with the given order.
func NewSorter(order SortOrder) *Sorter {
	return &Sorter{order: order}
}

// Apply returns a new slice of entries sorted by Date.
// Entries with a zero Date are placed at the end.
func (s *Sorter) Apply(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)

	sort.SliceStable(out, func(i, j int) bool {
		ai := out[i].Date
		bj := out[j].Date

		zeroTime := time.Time{}
		if ai == zeroTime && bj == zeroTime {
			return false
		}
		if ai == zeroTime {
			return false
		}
		if bj == zeroTime {
			return true
		}

		if s.order == SortAscending {
			return ai.Before(bj)
		}
		return ai.After(bj)
	})

	return out
}
