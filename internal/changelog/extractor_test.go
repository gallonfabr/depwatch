package changelog

import (
	"strings"
	"testing"
)

func TestNewExtractor_NotNil(t *testing.T) {
	e := NewExtractor()
	if e == nil {
		t.Fatal("expected non-nil Extractor")
	}
}

func TestExtractor_Apply_ExtractsIssueRef(t *testing.T) {
	e := NewExtractor()
	entries := []Entry{
		{Body: "Fixes #42 and closes GH-99"},
	}
	out := e.Apply(entries)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	hasTag := func(tag string) bool {
		for _, t := range out[0].Tags {
			if t == tag {
				return true
			}
		}
		return false
	}
	if !hasTag("issue:#42") && !hasTag("issue:42") {
		t.Errorf("expected issue tag for #42, got %v", out[0].Tags)
	}
}

func TestExtractor_Apply_ExtractsPRRef(t *testing.T) {
	e := NewExtractor()
	entries := []Entry{
		{Body: "Merged PR #7 into main"},
	}
	out := e.Apply(entries)
	found := false
	for _, tag := range out[0].Tags {
		if tag == "pr:7" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected pr:7 tag, got %v", out[0].Tags)
	}
}

func TestExtractor_Apply_ExtractsAuthor(t *testing.T) {
	e := NewExtractor()
	entries := []Entry{
		{Body: "Thanks @alice and @bob-dev for the patch"},
	}
	out := e.Apply(entries)
	authorTags := 0
	for _, tag := range out[0].Tags {
		if strings.HasPrefix(tag, "author:") {
			authorTags++
		}
	}
	if authorTags != 2 {
		t.Errorf("expected 2 author tags, got %d: %v", authorTags, out[0].Tags)
	}
}

func TestExtractor_Apply_NoRefs_EmptyTags(t *testing.T) {
	e := NewExtractor()
	entries := []Entry{
		{Body: "Minor internal refactor with no references"},
	}
	out := e.Apply(entries)
	if len(out[0].Tags) != 0 {
		t.Errorf("expected no tags, got %v", out[0].Tags)
	}
}

func TestExtractor_Apply_PreservesExistingTags(t *testing.T) {
	e := NewExtractor()
	entries := []Entry{
		{Body: "Fix #10", Tags: []string{"bugfix"}},
	}
	out := e.Apply(entries)
	hasBugfix := false
	for _, tag := range out[0].Tags {
		if tag == "bugfix" {
			hasBugfix = true
		}
	}
	if !hasBugfix {
		t.Errorf("expected existing 'bugfix' tag to be preserved, got %v", out[0].Tags)
	}
}

func TestExtractor_Apply_EmptyEntries(t *testing.T) {
	e := NewExtractor()
	out := e.Apply([]Entry{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
