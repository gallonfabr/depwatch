package changelog

import (
	"testing"
)

func sampleHighlightEntries() []Entry {
	return []Entry{
		{Version: "1.0.0", Body: "Fixed a critical security vulnerability in auth."},
		{Version: "1.1.0", Body: "Added new dashboard feature."},
		{Version: "1.2.0", Body: "Patch release: CVE-2024-1234 addressed."},
		{Version: "1.3.0", Body: "Performance improvements."},
	}
}

func TestNewHighlighter_NotNil(t *testing.T) {
	h := NewHighlighter()
	if h == nil {
		t.Fatal("expected non-nil Highlighter")
	}
}

func TestHighlighter_NoKeywords_NoHighlights(t *testing.T) {
	h := NewHighlighter()
	out := h.Apply(sampleHighlightEntries())
	for _, e := range out {
		if e.Meta[HighlightKey] == "true" {
			t.Errorf("entry %q should not be highlighted when no keywords set", e.Version)
		}
	}
}

func TestHighlighter_KeywordMatch_SetsFlag(t *testing.T) {
	h := NewHighlighter(WithHighlightKeywords("security"))
	out := h.Apply(sampleHighlightEntries())
	if out[0].Meta[HighlightKey] != "true" {
		t.Error("expected entry 0 to be highlighted (contains 'security')")
	}
	if out[1].Meta[HighlightKey] == "true" {
		t.Error("entry 1 should not be highlighted")
	}
}

func TestHighlighter_CaseInsensitiveMatch(t *testing.T) {
	h := NewHighlighter(WithHighlightKeywords("CVE"))
	out := h.Apply(sampleHighlightEntries())
	if out[2].Meta[HighlightKey] != "true" {
		t.Error("expected case-insensitive match on 'CVE'")
	}
}

func TestHighlighter_MultipleKeywords_AnyMatch(t *testing.T) {
	h := NewHighlighter(WithHighlightKeywords("security", "cve"))
	out := h.Apply(sampleHighlightEntries())
	if out[0].Meta[HighlightKey] != "true" {
		t.Error("entry 0 should match 'security'")
	}
	if out[2].Meta[HighlightKey] != "true" {
		t.Error("entry 2 should match 'cve'")
	}
	if out[3].Meta[HighlightKey] == "true" {
		t.Error("entry 3 should not be highlighted")
	}
}

func TestHighlighter_DoesNotMutateInput(t *testing.T) {
	original := sampleHighlightEntries()
	h := NewHighlighter(WithHighlightKeywords("security"))
	h.Apply(original)
	if original[0].Meta != nil && original[0].Meta[HighlightKey] == "true" {
		t.Error("Apply must not mutate the input slice entries")
	}
}

func TestHighlighter_EmptyEntries_ReturnsEmpty(t *testing.T) {
	h := NewHighlighter(WithHighlightKeywords("security"))
	out := h.Apply([]Entry{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
