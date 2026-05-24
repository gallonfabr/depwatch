package changelog

import (
	"regexp"
	"strings"
	"unicode"
)

// Normalizer cleans and standardizes changelog entry bodies
// before they are processed downstream.
type Normalizer struct {
	maxLength  int
	stripHTML   bool
	spaceRe    *regexp.Regexp
	htmlTagRe  *regexp.Regexp
}

// NormalizerOption configures a Normalizer.
type NormalizerOption func(*Normalizer)

// WithMaxLength sets the maximum body length (0 = unlimited).
func WithMaxLength(n int) NormalizerOption {
	return func(nr *Normalizer) { nr.maxLength = n }
}

// WithStripHTML enables removal of HTML tags from entry bodies.
func WithStripHTML(v bool) NormalizerOption {
	return func(nr *Normalizer) { nr.stripHTML = v }
}

// NewNormalizer constructs a Normalizer with optional configuration.
func NewNormalizer(opts ...NormalizerOption) *Normalizer {
	nr := &Normalizer{
		maxLength: 0,
		stripHTML:  false,
		spaceRe:   regexp.MustCompile(`[ \t]+`),
		htmlTagRe: regexp.MustCompile(`<[^>]+>`),
	}
	for _, o := range opts {
		o(nr)
	}
	return nr
}

// Normalize applies cleaning rules to a slice of Entry values and
// returns a new slice with updated bodies.
func (n *Normalizer) Normalize(entries []Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		e.Body = n.clean(e.Body)
		out = append(out, e)
	}
	return out
}

func (n *Normalizer) clean(s string) string {
	if n.stripHTML {
		s = n.htmlTagRe.ReplaceAllString(s, "")
	}
	// Collapse internal whitespace but preserve newlines.
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		l = n.spaceRe.ReplaceAllString(l, " ")
		lines[i] = strings.TrimRightFunc(l, unicode.IsSpace)
	}
	s = strings.Join(lines, "\n")
	s = strings.TrimSpace(s)
	if n.maxLength > 0 && len(s) > n.maxLength {
		s = s[:n.maxLength]
	}
	return s
}
