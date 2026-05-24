package changelog

import "strings"

// Tagger assigns free-form string tags to changelog entries based on
// configurable keyword rules. Unlike Labeler, tags are additive and
// multiple tags may be applied to a single entry.
type Tagger struct {
	rules map[string][]string // tag -> keywords that trigger it
}

// TaggerOption configures a Tagger.
type TaggerOption func(*Tagger)

// WithTagRule registers a tag that is applied whenever any of the
// provided keywords appear (case-insensitive) in an entry's body or
// version string.
func WithTagRule(tag string, keywords ...string) TaggerOption {
	return func(t *Tagger) {
		lower := make([]string, len(keywords))
		for i, k := range keywords {
			lower[i] = strings.ToLower(k)
		}
		t.rules[strings.ToLower(tag)] = append(t.rules[strings.ToLower(tag)], lower...)
	}
}

// NewTagger creates a Tagger with the supplied options.
func NewTagger(opts ...TaggerOption) *Tagger {
	t := &Tagger{rules: make(map[string][]string)}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Apply iterates over entries and appends matching tags to each entry's
// Tags slice. Existing tags are preserved.
func (t *Tagger) Apply(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)
	for i, e := range out {
		haystack := strings.ToLower(e.Body + " " + e.Version)
		for tag, keywords := range t.rules {
			if t.alreadyTagged(e.Tags, tag) {
				continue
			}
			for _, kw := range keywords {
				if strings.Contains(haystack, kw) {
					out[i].Tags = append(out[i].Tags, tag)
					break
				}
			}
		}
	}
	return out
}

func (t *Tagger) alreadyTagged(tags []string, tag string) bool {
	for _, existing := range tags {
		if strings.ToLower(existing) == tag {
			return true
		}
	}
	return false
}
