// Package changelog provides primitives for fetching, parsing, and
// processing dependency changelog entries.
//
// # Tagger
//
// Tagger assigns free-form string tags to [Entry] values based on
// keyword rules supplied at construction time via [WithTagRule].
//
// Unlike [Labeler], which maps entries to a fixed vocabulary of
// semantic labels, Tagger is open-ended: callers define both the tag
// names and the keywords that trigger them.
//
// Example:
//
//	tagger := changelog.NewTagger(
//		changelog.WithTagRule("security", "cve", "vulnerability"),
//		changelog.WithTagRule("beta",     "beta", "rc"),
//	)
//	tagged := tagger.Apply(entries)
//
// Tags are appended to [Entry].Tags; existing tags are preserved and
// duplicates are suppressed.
package changelog
