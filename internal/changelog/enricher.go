package changelog

import (
	"fmt"
	"strings"
)

// Enricher adds computed metadata to changelog entries.
type Enricher struct {
	baseURL string
}

// EnricherOption configures an Enricher.
type EnricherOption func(*Enricher)

// WithBaseURL sets a base URL used to build release links.
func WithBaseURL(u string) EnricherOption {
	return func(e *Enricher) {
		e.baseURL = strings.TrimRight(u, "/")
	}
}

// NewEnricher creates an Enricher with the given options.
func NewEnricher(opts ...EnricherOption) *Enricher {
	e := &Enricher{}
	for _, o := range opts {
		o(e)
	}
	return e
}

// Apply enriches a slice of entries in-place and returns them.
// It sets the Link field when a base URL is configured and
// ensures the Dependency field is non-empty.
func (e *Enricher) Apply(dep string, entries []Entry) []Entry {
	for i := range entries {
		if entries[i].Dependency == "" {
			entries[i].Dependency = dep
		}
		if e.baseURL != "" && entries[i].Link == "" && entries[i].Version != "" {
			entries[i].Link = fmt.Sprintf("%s/%s", e.baseURL, entries[i].Version)
		}
	}
	return entries
}
