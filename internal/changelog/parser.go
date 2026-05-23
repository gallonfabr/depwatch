package changelog

import (
	"regexp"
	"strings"
	"time"
)

// Entry represents a single changelog entry for a version.
type Entry struct {
	Version   string
	Date      time.Time
	Body      string
	RawDate   string
}

// Parser parses raw changelog text into structured entries.
type Parser struct {
	versionRe *regexp.Regexp
}

// NewParser creates a Parser that recognises common changelog headings.
// Supported heading formats:
//   ## [1.2.3] - 2024-01-15
//   ## 1.2.3 (2024-01-15)
//   ## v1.2.3 — 2024-01-15
func NewParser() *Parser {
	return &Parser{
		versionRe: regexp.MustCompile(
			`(?i)^#{1,3}\s+v?(\d+\.\d+\.\d+[^\s]*)\s*[\-\(—]?\s*(\d{4}-\d{2}-\d{2})?`,
		),
	}
}

// Parse splits raw changelog text into a slice of Entry values.
// Entries are returned in the order they appear in the text.
func (p *Parser) Parse(raw string) []Entry {
	lines := strings.Split(raw, "\n")
	var entries []Entry
	var current *Entry
	var bodyLines []string

	flush := func() {
		if current != nil {
			current.Body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
			entries = append(entries, *current)
			current = nil
			bodyLines = nil
		}
	}

	for _, line := range lines {
		matches := p.versionRe.FindStringSubmatch(line)
		if len(matches) >= 2 {
			flush()
			entry := &Entry{Version: matches[1]}
			if len(matches) >= 3 && matches[2] != "" {
				entry.RawDate = matches[2]
				if t, err := time.Parse("2006-01-02", matches[2]); err == nil {
					entry.Date = t
				}
			}
			current = entry
		} else if current != nil {
			bodyLines = append(bodyLines, line)
		}
	}
	flush()
	return entries
}
