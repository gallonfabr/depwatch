package changelog

import "strings"

// Scorer assigns a relevance score to a changelog Entry based on configurable
// keywords. Higher scores indicate higher relevance.
type Scorer struct {
	keywords []string
}

// ScorerOption configures a Scorer.
type ScorerOption func(*Scorer)

// WithKeywords adds keywords that boost an entry's score when found in its
// body or version string.
func WithKeywords(kw ...string) ScorerOption {
	return func(s *Scorer) {
		s.keywords = append(s.keywords, kw...)
	}
}

// NewScorer creates a Scorer with the supplied options.
func NewScorer(opts ...ScorerOption) *Scorer {
	s := &Scorer{}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Score returns a non-negative integer representing how relevant the entry is.
// Each keyword match in the body (case-insensitive) contributes 1 point.
// A security-related keyword contributes 2 points.
func (s *Scorer) Score(e Entry) int {
	score := 0
	body := strings.ToLower(e.Body)
	version := strings.ToLower(e.Version.String())

	for _, kw := range s.keywords {
		lower := strings.ToLower(kw)
		count := strings.Count(body, lower) + strings.Count(version, lower)
		if count == 0 {
			continue
		}
		weight := 1
		if isSecurity(lower) {
			weight = 2
		}
		score += count * weight
	}
	return score
}

// isSecurity returns true for keywords that indicate a security-relevant change.
func isSecurity(kw string) bool {
	securityTerms := []string{"security", "cve", "vuln", "exploit", "patch", "advisory"}
	for _, t := range securityTerms {
		if strings.Contains(kw, t) {
			return true
		}
	}
	return false
}
