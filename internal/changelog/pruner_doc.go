// Package changelog provides Pruner, which removes changelog entries that
// exceed a configurable maximum age.
//
// # Usage
//
//	pruner := changelog.NewPruner(
//		changelog.WithMaxAge(14 * 24 * time.Hour),
//	)
//	filtered := pruner.Apply(entries)
//
// Entries with a zero date are always retained so that entries without
// reliable date information are never silently discarded.
package changelog
