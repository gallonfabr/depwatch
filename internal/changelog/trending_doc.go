// Package changelog provides primitives for fetching, parsing, and
// processing dependency changelogs.
//
// Trending
//
// The Trending analyser aggregates scored entries by dependency name and
// returns a ranked list of the most active or impactful dependencies over a
// given set of changelog entries.
//
// Usage:
//
//	tr := changelog.NewTrending(5) // top 5 dependencies
//	result := tr.Analyse(entries)
//	for _, e := range result {
//		fmt.Printf("%s: score=%d count=%d\n", e.Dependency, e.Score, e.Count)
//	}
package changelog
