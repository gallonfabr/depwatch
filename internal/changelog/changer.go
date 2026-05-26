package changelog

import "strings"

// Changer detects whether a new set of entries represents a meaningful change
// compared to a previously seen snapshot. It is useful for suppressing
// notifications when nothing has actually changed between polling cycles.
type Changer struct {
	store *SnapshotStore
}

// NewChanger returns a Changer backed by the given SnapshotStore.
// It panics if store is nil.
func NewChanger(store *SnapshotStore) *Changer {
	if store == nil {
		panic("changer: store must not be nil")
	}
	return &Changer{store: store}
}

// HasChanged reports whether entries for the given dependency differ from the
// last saved snapshot. It returns true when there is no previous snapshot.
func (c *Changer) HasChanged(dep string, entries []Entry) bool {
	prev, ok := c.store.Get(dep)
	if !ok {
		return len(entries) > 0
	}
	if len(prev) != len(entries) {
		return true
	}
	for i, e := range entries {
		p := prev[i]
		if !strings.EqualFold(e.Version, p.Version) || e.Date != p.Date {
			return true
		}
	}
	return false
}

// Commit saves the current entries as the new snapshot for dep.
func (c *Changer) Commit(dep string, entries []Entry) {
	c.store.Save(dep, entries)
}
