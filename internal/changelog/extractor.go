package changelog

import (
	"regexp"
	"strings"
)

// Extractor pulls structured metadata (issue refs, PR links, authors) from
// an Entry's body and stores them in the entry's Tags slice or a dedicated
// Metadata map if one is available via the Entry type.
//
// It operates as a Transformer so it can be composed in a Pipeline or Chain.
type Extractor struct {
	issueRe  *regexp.Regexp
	prRe     *regexp.Regexp
	authorRe *regexp.Regexp
}

// ExtractorOption configures an Extractor.
type ExtractorOption func(*Extractor)

// NewExtractor returns an Extractor with sensible defaults.
// Patterns match common conventions:
//   - Issues:  #123 or GH-123
//   - PRs:     PR #123 or pr-123
//   - Authors: @username
func NewExtractor(opts ...ExtractorOption) *Extractor {
	e := &Extractor{
		issueRe:  regexp.MustCompile(`(?i)(?:GH-|#)(\d+)`),
		prRe:     regexp.MustCompile(`(?i)(?:pr[\s#-]+(\d+))`),
		authorRe: regexp.MustCompile(`@([A-Za-z0-9_-]+)`),
	}
	for _, o := range opts {
		o(e)
	}
	return e
}

// Apply implements Transformer. It annotates each Entry with tags derived
// from extracted references found in the body text.
func (e *Extractor) Apply(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	for i, en := range entries {
		out[i] = e.annotate(en)
	}
	return out
}

func (e *Extractor) annotate(en Entry) Entry {
	body := en.Body
	tagSet := make(map[string]struct{})
	for _, t := range en.Tags {
		tagSet[t] = struct{}{}
	}

	for _, m := range e.issueRe.FindAllString(body, -1) {
		tag := "issue:" + strings.TrimLeft(strings.ToUpper(m), "GHgh- ")
		tagSet[tag] = struct{}{}
	}
	for _, m := range e.prRe.FindAllStringSubmatch(body, -1) {
		if len(m) > 1 {
			tagSet["pr:"+m[1]] = struct{}{}
		}
	}
	for _, m := range e.authorRe.FindAllStringSubmatch(body, -1) {
		if len(m) > 1 {
			tagSet["author:"+m[1]] = struct{}{}
		}
	}

	tags := make([]string, 0, len(tagSet))
	for t := range tagSet {
		tags = append(tags, t)
	}
	en.Tags = tags
	return en
}
