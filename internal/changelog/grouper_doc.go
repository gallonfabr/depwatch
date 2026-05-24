// Package changelog provides utilities for fetching, parsing, and
// processing dependency changelogs.
//
// Grouper partitions a flat list of [Entry] values into named [Group]
// buckets based on the first label attached to each entry (see Labeler).
//
// Example usage:
//
//	g := changelog.NewGrouper(
//		changelog.WithGroupOrder("security", "feature", "bugfix"),
//		changelog.WithFallbackLabel("other"),
//	)
//	groups := g.Apply(entries)
//	for _, grp := range groups {
//		fmt.Println(grp.Label, len(grp.Entries))
//	}
package changelog
