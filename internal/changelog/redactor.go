package changelog

import (
	"regexp"
	"strings"
)

// Redactor replaces sensitive patterns in entry bodies with a placeholder.
// Useful for stripping tokens, API keys, or internal URLs before delivery.
type Redactor struct {
	patterns    []*regexp.Regexp
	placeholder string
}

// RedactorOption configures a Redactor.
type RedactorOption func(*Redactor)

// WithRedactPattern adds a compiled regular expression whose matches will be
// replaced by the placeholder string.
func WithRedactPattern(pattern string) RedactorOption {
	return func(r *Redactor) {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return
		}
		r.patterns = append(r.patterns, re)
	}
}

// WithRedactPlaceholder overrides the default replacement string ("[REDACTED]").
func WithRedactPlaceholder(s string) RedactorOption {
	return func(r *Redactor) {
		if s != "" {
			r.placeholder = s
		}
	}
}

// NewRedactor creates a Redactor with the supplied options.
func NewRedactor(opts ...RedactorOption) *Redactor {
	r := &Redactor{placeholder: "[REDACTED]"}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Apply iterates over entries and redacts sensitive content from Body and Link.
func (r *Redactor) Apply(entries []Entry) []Entry {
	if len(r.patterns) == 0 {
		return entries
	}
	out := make([]Entry, len(entries))
	for i, e := range entries {
		e.Body = r.redact(e.Body)
		e.Link = r.redact(e.Link)
		out[i] = e
	}
	return out
}

func (r *Redactor) redact(s string) string {
	for _, re := range r.patterns {
		s = re.ReplaceAllString(s, r.placeholder)
	}
	return strings.TrimSpace(s)
}
