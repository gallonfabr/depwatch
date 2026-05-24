package changelog

import (
	"testing"
	"time"
)

func sampleEnrichEntries() []Entry {
	return []Entry{
		{Version: "v1.2.0", Date: time.Now(), Body: "fix: something"},
		{Version: "v1.1.0", Date: time.Now(), Body: "feat: other"},
	}
}

func TestNewEnricher_NotNil(t *testing.T) {
	e := NewEnricher()
	if e == nil {
		t.Fatal("expected non-nil Enricher")
	}
}

func TestEnricher_SetsDependency(t *testing.T) {
	e := NewEnricher()
	entries := sampleEnrichEntries()
	result := e.Apply("mylib", entries)
	for _, en := range result {
		if en.Dependency != "mylib" {
			t.Errorf("expected dependency 'mylib', got %q", en.Dependency)
		}
	}
}

func TestEnricher_DoesNotOverwriteDependency(t *testing.T) {
	e := NewEnricher()
	entries := []Entry{{Version: "v1.0.0", Dependency: "existing"}}
	result := e.Apply("newdep", entries)
	if result[0].Dependency != "existing" {
		t.Errorf("expected 'existing', got %q", result[0].Dependency)
	}
}

func TestEnricher_SetsLinkWithBaseURL(t *testing.T) {
	e := NewEnricher(WithBaseURL("https://github.com/org/repo/releases/tag"))
	entries := sampleEnrichEntries()
	result := e.Apply("repo", entries)
	for _, en := range result {
		expected := "https://github.com/org/repo/releases/tag/" + en.Version
		if en.Link != expected {
			t.Errorf("expected link %q, got %q", expected, en.Link)
		}
	}
}

func TestEnricher_NoBaseURL_LinkEmpty(t *testing.T) {
	e := NewEnricher()
	entries := sampleEnrichEntries()
	result := e.Apply("repo", entries)
	for _, en := range result {
		if en.Link != "" {
			t.Errorf("expected empty link, got %q", en.Link)
		}
	}
}

func TestEnricher_DoesNotOverwriteExistingLink(t *testing.T) {
	e := NewEnricher(WithBaseURL("https://example.com"))
	entries := []Entry{{Version: "v2.0.0", Link: "https://custom.link/v2.0.0"}}
	result := e.Apply("dep", entries)
	if result[0].Link != "https://custom.link/v2.0.0" {
		t.Errorf("expected original link to be preserved, got %q", result[0].Link)
	}
}

func TestEnricher_EmptyEntries(t *testing.T) {
	e := NewEnricher(WithBaseURL("https://example.com"))
	result := e.Apply("dep", []Entry{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestWithBaseURL_TrimsTrailingSlash(t *testing.T) {
	e := NewEnricher(WithBaseURL("https://example.com/releases/"))
	entries := []Entry{{Version: "v1.0.0"}}
	result := e.Apply("dep", entries)
	if result[0].Link != "https://example.com/releases/v1.0.0" {
		t.Errorf("unexpected link: %q", result[0].Link)
	}
}
