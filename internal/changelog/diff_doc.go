// Package changelog provides utilities for fetching, parsing, and processing
// dependency changelogs.
//
// # Diff
//
// Diff compares two snapshots of changelog entries and isolates entries that
// are new in the current snapshot.  It is used by the watcher to avoid
// re-reporting already-seen releases.
//
// Usage:
//
//	d := changelog.NewDiff()
//	novel := d.Apply(previousEntries, currentEntries)
//	summary := d.Summarise(novel)
package changelog
