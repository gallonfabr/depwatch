package changelog

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time capture of fetched entries for a dependency.
type Snapshot struct {
	Dependency string
	Entries    []Entry
	CapturedAt time.Time
}

// SnapshotStore keeps the most recent snapshot per dependency.
type SnapshotStore struct {
	mu    sync.RWMutex
	store map[string]Snapshot
}

// NewSnapshotStore creates an empty SnapshotStore.
func NewSnapshotStore() *SnapshotStore {
	return &SnapshotStore{
		store: make(map[string]Snapshot),
	}
}

// Save records a snapshot for the given dependency, replacing any previous one.
func (s *SnapshotStore) Save(dep string, entries []Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	copy := make([]Entry, len(entries))
	for i, e := range entries {
		copy[i] = e
	}

	s.store[dep] = Snapshot{
		Dependency: dep,
		Entries:    copy,
		CapturedAt: time.Now(),
	}
}

// Get returns the latest snapshot for a dependency and whether it exists.
func (s *SnapshotStore) Get(dep string) (Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap, ok := s.store[dep]
	return snap, ok
}

// All returns a copy of all stored snapshots.
func (s *SnapshotStore) All() []Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Snapshot, 0, len(s.store))
	for _, snap := range s.store {
		result = append(result, snap)
	}
	return result
}

// Clear removes the snapshot for the given dependency.
func (s *SnapshotStore) Clear(dep string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, dep)
}
