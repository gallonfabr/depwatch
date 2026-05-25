// Package changelog provides primitives for fetching, parsing, and
// processing dependency changelogs.
//
// SemVerFilter
//
// SemVerFilter is a pipeline stage that keeps only those changelog
// entries whose semantic version falls within an inclusive [min, max]
// range.  Either bound may be omitted.
//
// Example usage:
//
//	f := changelog.NewSemVerFilter(
//		changelog.WithMinVersion("v1.2.0"),
//		changelog.WithMaxVersion("v2.0.0"),
//	)
//	filtered := f.Apply(entries)
package changelog
