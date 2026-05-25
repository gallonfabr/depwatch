// Package changelog provides primitives for fetching, parsing, and
// processing dependency changelog entries.
//
// # Suppressor
//
// Suppressor is a transformer that removes changelog entries belonging to
// explicitly suppressed dependencies. This is useful when certain packages
// produce high-volume or low-signal releases that should be excluded from
// digests without removing them from configuration entirely.
//
// Example usage:
//
//	s := changelog.NewSuppressor([]string{"some-noisy-dep", "another-dep"})
//	filtered := s.Apply(entries)
package changelog
