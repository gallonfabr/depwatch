// Package changelog provides primitives for fetching, parsing, and
// processing dependency changelogs.
//
// Sanitizer cleans raw entry bodies fetched from external sources.
// It removes ASCII control characters (while preserving tabs and newlines),
// normalises line endings to LF, and optionally truncates bodies to a
// configurable maximum number of Unicode code points.
//
// Usage:
//
//	s := changelog.NewSanitizer(changelog.WithMaxRunes(500))
//	clean := s.Apply(entries)
package changelog
