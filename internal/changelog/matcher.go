package changelog

import (
	"regexp"
	"strings"
)

// MatchRule defines a single pattern-based rule for matching changelog entries.
type MatchRule struct {
	Pattern  *regexp.Regexp
	Field    string // "title", "body", or "version"
}

// Matcher filters entries that satisfy at least one match rule.
type Matcher struct {
	rules []MatchRule
}

// MatcherOption configures a Matcher.
type MatcherOption func(*Matcher)

// WithMatchRule adds a compiled regex rule targeting a specific field.
func WithMatchRule(field, pattern string) MatcherOption {
	re := regexp.MustCompile(pattern)
	return func(m *Matcher) {
		m.rules = append(m.rules, MatchRule{Pattern: re, Field: field})
	}
}

// NewMatcher constructs a Matcher with the given options.
func NewMatcher(opts ...MatcherOption) *Matcher {
	m := &Matcher{}
	for _, o := range opts {
		o(m)
	}
	return m
}

// Apply returns only entries that match at least one rule.
// If no rules are configured, all entries are returned unchanged.
func (m *Matcher) Apply(entries []Entry) []Entry {
	if len(m.rules) == 0 {
		return entries
	}
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if m.matches(e) {
			out = append(out, e)
		}
	}
	return out
}

func (m *Matcher) matches(e Entry) bool {
	for _, r := range m.rules {
		var target string
		switch strings.ToLower(r.Field) {
		case "title":
			target = e.Version
		case "body":
			target = e.Body
		case "version":
			target = e.Version
		default:
			target = e.Body
		}
		if r.Pattern.MatchString(target) {
			return true
		}
	}
	return false
}
