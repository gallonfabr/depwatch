// Package changelog provides primitives for fetching, parsing, and processing
// dependency changelogs.
//
// # Scorer and Ranker
//
// Scorer assigns a numeric relevance score to a changelog Entry based on
// user-supplied keywords. Security-related keywords ("security", "cve",
// "vuln", etc.) receive double weight.
//
// Ranker wraps a Scorer and returns entries sorted from highest to lowest
// score. Use Ranker.TopN to retrieve only the most relevant entries, which is
// useful when building digest summaries that should highlight critical changes.
//
// Example:
//
//	scorer := changelog.NewScorer(
//		changelog.WithKeywords("security", "breaking", "fix"),
//	)
//	ranker := changelog.NewRanker(scorer)
//	top5 := ranker.TopN(entries, 5)
package changelog
