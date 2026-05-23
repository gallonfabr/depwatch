// Package store provides a simple persistent key-value store for tracking
// the last-seen changelog version per dependency.
package store

import (
	"encoding/json"
	"os"
	"sync"
)

// Store persists the last-seen version for each dependency so that depwatch
// can emit only new changelog entries on subsequent polls.
type Store struct {
	mu   sync.RWMutex
	path string
	data map[string]string
}

// New opens (or creates) a JSON file at path and returns a ready Store.
func New(path string) (*Store, error) {
	s := &Store{
		path: path,
		data: make(map[string]string),
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(bytes, &s.data); err != nil {
		return nil, err
	}
	return s, nil
}

// LastSeen returns the last-seen version string for the given dependency key.
// It returns an empty string when the key has never been recorded.
func (s *Store) LastSeen(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[key]
}

// SetLastSeen records version as the latest seen version for key and
// immediately flushes the store to disk.
func (s *Store) SetLastSeen(key, version string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = version
	return s.flush()
}

// flush writes the current in-memory state to the JSON file.
// Callers must hold s.mu (write lock).
func (s *Store) flush() error {
	bytes, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, bytes, 0o644)
}
