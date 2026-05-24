package changelog

import "sort"

// Ranker sorts a slice of entries by relevance score (descending) using a
// Scorer. Entries with equal scores preserve their original relative order
// (stable sort).
type Ranker struct {
	scorer *Scorer
}

// NewRanker creates a Ranker backed by the provided Scorer.
// If scorer is nil, NewRanker panics.
func NewRanker(scorer *Scorer) *Ranker {
	if scorer == nil {
		panic("changelog: NewRanker requires a non-nil Scorer")
	}
	return &Ranker{scorer: scorer}
}

// Rank returns a new slice of entries ordered from highest to lowest score.
// The original slice is not modified.
func (r *Ranker) Rank(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)

	scores := make([]int, len(out))
	for i, e := range out {
		scores[i] = r.scorer.Score(e)
	}

	// Stable sort so equal-score entries keep their original order.
	sort.SliceStable(out, func(i, j int) bool {
		return scores[i] > scores[j]
	})

	return out
}

// TopN returns at most n highest-ranked entries.
// If n <= 0 all entries are returned.
func (r *Ranker) TopN(entries []Entry, n int) []Entry {
	ranked := r.Rank(entries)
	if n <= 0 || n >= len(ranked) {
		return ranked
	}
	return ranked[:n]
}
