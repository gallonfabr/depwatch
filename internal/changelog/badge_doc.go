// Package changelog provides primitives for fetching, parsing, and
// processing dependency changelogs.
//
// # Badger
//
// Badger attaches coloured Badge values to changelog entries based on
// their existing Tags. Rules are registered with WithBadgeRule and
// evaluated in order; duplicate badges (same Label) are silently dropped.
//
// Example:
//
//	b := changelog.NewBadger(
//		changelog.WithBadgeRule("security", "Security", "#e11d48"),
//		changelog.WithBadgeRule("breaking", "Breaking", "#f97316"),
//	)
//	tagged := b.Apply(entries)
package changelog
