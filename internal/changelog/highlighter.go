package changelog

import (
	"strings"
)

// Highlighter marks entries whose body or version string contains any of the
// configured keywords by setting a boolean flag in the entry's metadata.
// It is intended to surface high-priority items (e.g. security fixes) in
// digest renderers without altering the underlying content.

// HighlightKey is the metadata key written by Highlighter.
const HighlightKey = "highlight"

// Highlighter is a Transformer that flags matching entries.
type Highlighter struct {
	keywords []string
}

// HighlighterOption configures a Highlighter.
type HighlighterOption func(*Highlighter)

// WithHighlightKeywords adds keywords that trigger highlighting.
func WithHighlightKeywords(kw ...string) HighlighterOption {
	return func(h *Highlighter) {
		for _, k := range kw {
			h.keywords = append(h.keywords, strings.ToLower(k))
		}
	}
}

// NewHighlighter returns a Highlighter configured with the supplied options.
// At least one keyword must be provided; if none are given the transformer is
// a no-op (no entries are highlighted).
func NewHighlighter(opts ...HighlighterOption) *Highlighter {
	h := &Highlighter{}
	for _, o := range opts {
		o(h)
	}
	return h
}

// Apply implements Transformer. It sets Entry.Meta[HighlightKey] = "true" for
// every entry whose lowercased body or version contains a configured keyword.
func (h *Highlighter) Apply(entries []Entry) []Entry {
	if len(h.keywords) == 0 {
		return entries
	}
	out := make([]Entry, len(entries))
	for i, e := range entries {
		haystack := strings.ToLower(e.Body + " " + e.Version)
		for _, kw := range h.keywords {
			if strings.Contains(haystack, kw) {
				if e.Meta == nil {
					e.Meta = map[string]string{}
				}
				e.Meta[HighlightKey] = "true"
				break
			}
		}
		out[i] = e
	}
	return out
}
