package changelog

import "sort"

// TrendingEntry pairs a dependency name with its aggregated score.
type TrendingEntry struct {
	Dependency string
	Score      int
	Count      int
}

// Trending analyses a slice of entries and returns a ranked list of
// dependencies ordered by total score descending, then by entry count.
type Trending struct {
	topN int
}

// NewTrending creates a Trending analyser. topN controls how many
// dependencies are returned; pass 0 to return all.
func NewTrending(topN int) *Trending {
	return &Trending{topN: topN}
}

// Analyse aggregates scores from entries and returns ranked TrendingEntry
// values.
func (t *Trending) Analyse(entries []Entry) []TrendingEntry {
	type agg struct {
		score int
		count int
	}

	accum := make(map[string]*agg)
	for _, e := range entries {
		key := e.Dependency
		if key == "" {
			continue
		}
		a, ok := accum[key]
		if !ok {
			a = &agg{}
			accum[key] = a
		}
		a.score += e.Score
		a.count++
	}

	result := make([]TrendingEntry, 0, len(accum))
	for dep, a := range accum {
		result = append(result, TrendingEntry{
			Dependency: dep,
			Score:      a.score,
			Count:      a.count,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Score != result[j].Score {
			return result[i].Score > result[j].Score
		}
		return result[i].Count > result[j].Count
	})

	if t.topN > 0 && len(result) > t.topN {
		result = result[:t.topN]
	}
	return result
}
