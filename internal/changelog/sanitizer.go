package changelog

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Sanitizer cleans entry bodies by removing control characters,
// normalising line endings, and optionally truncating to a max rune count.
type Sanitizer struct {
	maxRunes int
	controlRe *regexp.Regexp
}

// SanitizerOption configures a Sanitizer.
type SanitizerOption func(*Sanitizer)

// WithMaxRunes sets the maximum number of Unicode code points kept per body.
func WithMaxRunes(n int) SanitizerOption {
	return func(s *Sanitizer) {
		if n > 0 {
			s.maxRunes = n
		}
	}
}

// NewSanitizer returns a Sanitizer ready to use.
func NewSanitizer(opts ...SanitizerOption) *Sanitizer {
	s := &Sanitizer{
		controlRe: regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`),
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Apply implements Transformer.
func (s *Sanitizer) Apply(entries []Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		e.Body = s.clean(e.Body)
		out = append(out, e)
	}
	return out
}

func (s *Sanitizer) clean(body string) string {
	// Normalise Windows line endings.
	body = strings.ReplaceAll(body, "\r\n", "\n")
	body = strings.ReplaceAll(body, "\r", "\n")

	// Strip control characters (keep \t and \n).
	body = s.controlRe.ReplaceAllString(body, "")

	// Truncate if a max rune limit is set.
	if s.maxRunes > 0 && utf8.RuneCountInString(body) > s.maxRunes {
		runes := []rune(body)
		body = string(runes[:s.maxRunes])
	}

	return strings.TrimSpace(body)
}
