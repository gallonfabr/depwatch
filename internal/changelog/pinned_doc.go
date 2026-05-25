// Package changelog — pinned.go
//
// Pinner provides a lightweight mechanism to suppress digest notifications
// for specific (dependency, version) pairs that an operator has explicitly
// acknowledged or wishes to ignore temporarily.
//
// Usage:
//
//	pins := []changelog.PinnedEntry{
//		{Dependency: "react", Version: "18.2.0", Reason: "already reviewed"},
//	}
//	pinner := changelog.NewPinner(pins)
//	filtered := pinner.Apply(entries)
//
// Pinned entries are stored in memory for the lifetime of the Pinner.
// To persist pins across restarts, serialize []PinnedEntry to the config
// or an external store and reload on startup.
package changelog
