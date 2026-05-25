package changelog

import "sort"

// AggregateStats holds summary statistics for a set of changelog entries.
type AggregateStats struct {
	TotalEntries    int
	ByDependency    map[string]int
	ByLabel         map[string]int
	HighlightedCount int
}

// Aggregator computes summary statistics over a slice of Entry values.
type Aggregator struct{}

// NewAggregator returns a new Aggregator.
func NewAggregator() *Aggregator {
	return &Aggregator{}
}

// Aggregate computes statistics from the provided entries.
func (a *Aggregator) Aggregate(entries []Entry) AggregateStats {
	stats := AggregateStats{
		TotalEntries: len(entries),
		ByDependency: make(map[string]int),
		ByLabel:      make(map[string]int),
	}

	for _, e := range entries {
		if e.Dependency != "" {
			stats.ByDependency[e.Dependency]++
		}
		for _, lbl := range e.Labels {
			stats.ByLabel[lbl]++
		}
		if e.Highlighted {
			stats.HighlightedCount++
		}
	}

	return stats
}

// TopDependencies returns dependency names sorted by entry count descending,
// capped at n. If n <= 0 all dependencies are returned.
func (a *Aggregator) TopDependencies(stats AggregateStats, n int) []string {
	type kv struct {
		Key   string
		Count int
	}
	pairs := make([]kv, 0, len(stats.ByDependency))
	for k, v := range stats.ByDependency {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count != pairs[j].Count {
			return pairs[i].Count > pairs[j].Count
		}
		return pairs[i].Key < pairs[j].Key
	})

	out := make([]string, 0, len(pairs))
	for _, p := range pairs {
		out = append(out, p.Key)
	}
	if n > 0 && n < len(out) {
		return out[:n]
	}
	return out
}
