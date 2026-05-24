package changelog

import (
	"strings"
	"unicode"
)

// Summarizer truncates each entry's body to a short summary suitable for
// digest previews. It trims whitespace, strips newlines, and caps the output
// at a configurable rune limit, appending an ellipsis when truncated.
type Summarizer struct {
	maxRunes int
}

// SummarizerOption configures a Summarizer.
type SummarizerOption func(*Summarizer)

// WithSummaryLength sets the maximum number of runes kept per summary.
// Values below 1 are ignored.
func WithSummaryLength(n int) SummarizerOption {
	return func(s *Summarizer) {
		if n > 0 {
			s.maxRunes = n
		}
	}
}

const defaultSummaryLength = 120

// NewSummarizer returns a Summarizer with the given options applied.
func NewSummarizer(opts ...SummarizerOption) *Summarizer {
	s := &Summarizer{maxRunes: defaultSummaryLength}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Apply implements Transformer. It rewrites each entry's Body with a compact
// single-line summary.
func (s *Summarizer) Apply(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		e.Body = s.summarize(e.Body)
		out[i] = e
	}
	return out
}

func (s *Summarizer) summarize(body string) string {
	// Collapse all whitespace sequences (including newlines) to a single space.
	fields := strings.FieldsFunc(body, func(r rune) bool {
		return unicode.IsSpace(r)
	})
	flat := strings.Join(fields, " ")

	runes := []rune(flat)
	if len(runes) <= s.maxRunes {
		return flat
	}
	return string(runes[:s.maxRunes]) + "…"
}
