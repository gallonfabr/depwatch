package digest

import (
	"fmt"
	"strings"
	"time"

	"github.com/yourorg/depwatch/internal/changelog"
)

// Entry represents a single dependency update in a digest.
type Entry struct {
	Dependency string
	Version    string
	Date       time.Time
	Body       string
}

// Digest holds a collection of dependency update entries.
type Digest struct {
	GeneratedAt time.Time
	Entries     []Entry
}

// Builder constructs a Digest from parsed changelog entries.
type Builder struct{}

// NewBuilder returns a new Builder instance.
func NewBuilder() *Builder {
	return &Builder{}
}

// Build creates a Digest from a map of dependency name to parsed changelog entries.
func (b *Builder) Build(updates map[string][]changelog.Entry) *Digest {
	d := &Digest{
		GeneratedAt: time.Now().UTC(),
	}
	for dep, entries := range updates {
		for _, e := range entries {
			d.Entries = append(d.Entries, Entry{
				Dependency: dep,
				Version:    e.Version,
				Date:       e.Date,
				Body:       e.Body,
			})
		}
	}
	return d
}

// FormatText renders the digest as a plain-text string.
func (d *Digest) FormatText() string {
	if len(d.Entries) == 0 {
		return "No dependency updates found."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Dependency Digest — %s\n", d.GeneratedAt.Format("2006-01-02 15:04 UTC")))
	sb.WriteString(strings.Repeat("=", 48) + "\n")
	for _, e := range d.Entries {
		sb.WriteString(fmt.Sprintf("[%s] %s (%s)\n", e.Dependency, e.Version, e.Date.Format("2006-01-02")))
		if e.Body != "" {
			sb.WriteString(e.Body + "\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
