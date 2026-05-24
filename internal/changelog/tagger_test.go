package changelog

import (
	"testing"
	"time"
)

var sampleTaggerEntries = []Entry{
	{Dependency: "libA", Version: "1.0.0", Date: time.Now(), Body: "Fixed a critical security vulnerability in auth module"},
	{Dependency: "libB", Version: "2.1.0", Date: time.Now(), Body: "Added new dashboard feature and improved UX"},
	{Dependency: "libC", Version: "3.0.0-beta", Date: time.Now(), Body: "Performance improvements and bug fixes"},
	{Dependency: "libD", Version: "1.2.3", Date: time.Now(), Body: "Routine maintenance release"},
}

func TestNewTagger_NotNil(t *testing.T) {
	tagger := NewTagger()
	if tagger == nil {
		t.Fatal("expected non-nil Tagger")
	}
}

func TestNewTagger_NoRules_NoTags(t *testing.T) {
	tagger := NewTagger()
	out := tagger.Apply(sampleTaggerEntries)
	for _, e := range out {
		if len(e.Tags) != 0 {
			t.Errorf("expected no tags for %s, got %v", e.Dependency, e.Tags)
		}
	}
}

func TestTagger_SecurityKeyword(t *testing.T) {
	tagger := NewTagger(WithTagRule("security", "vulnerability", "cve"))
	out := tagger.Apply(sampleTaggerEntries)

	if !containsTag(out[0].Tags, "security") {
		t.Errorf("libA should have 'security' tag, got %v", out[0].Tags)
	}
	for _, e := range out[1:] {
		if containsTag(e.Tags, "security") {
			t.Errorf("%s should not have 'security' tag", e.Dependency)
		}
	}
}

func TestTagger_MultipleRules(t *testing.T) {
	tagger := NewTagger(
		WithTagRule("security", "security"),
		WithTagRule("feature", "feature", "dashboard"),
		WithTagRule("beta", "beta"),
	)
	out := tagger.Apply(sampleTaggerEntries)

	if !containsTag(out[0].Tags, "security") {
		t.Error("libA: expected 'security'")
	}
	if !containsTag(out[1].Tags, "feature") {
		t.Error("libB: expected 'feature'")
	}
	if !containsTag(out[2].Tags, "beta") {
		t.Error("libC: expected 'beta'")
	}
}

func TestTagger_NoDuplicateTags(t *testing.T) {
	entry := Entry{Dependency: "libX", Version: "1.0.0", Body: "security fix", Tags: []string{"security"}}
	tagger := NewTagger(WithTagRule("security", "security"))
	out := tagger.Apply([]Entry{entry})

	count := 0
	for _, tag := range out[0].Tags {
		if tag == "security" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly 1 'security' tag, got %d", count)
	}
}

func TestTagger_CaseInsensitiveMatch(t *testing.T) {
	entry := Entry{Dependency: "libY", Version: "1.0.0", Body: "BREAKING CHANGE detected"}
	tagger := NewTagger(WithTagRule("breaking", "breaking change"))
	out := tagger.Apply([]Entry{entry})

	if !containsTag(out[0].Tags, "breaking") {
		t.Errorf("expected 'breaking' tag, got %v", out[0].Tags)
	}
}

func TestTagger_DoesNotMutateInput(t *testing.T) {
	original := []Entry{{Dependency: "libZ", Version: "1.0.0", Body: "security patch"}}
	tagger := NewTagger(WithTagRule("security", "security"))
	_ = tagger.Apply(original)

	if len(original[0].Tags) != 0 {
		t.Error("Apply must not mutate the original slice")
	}
}

func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
