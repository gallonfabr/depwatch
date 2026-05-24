// Package changelog provides types and utilities for fetching,
// parsing, and processing dependency changelog entries.
package changelog

import "time"

// Entry represents a single versioned release in a dependency changelog.
type Entry struct {
	// Dependency is the name of the dependency this entry belongs to.
	Dependency string

	// Version is the release version string, e.g. "v1.2.3".
	Version string

	// Date is the release date of this entry.
	Date time.Time

	// Body contains the human-readable changelog text for this release.
	Body string

	// Link is an optional URL pointing to the release page or tag.
	Link string
}
