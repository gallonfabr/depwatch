package changelog

import (
	"strings"
	"testing"
	"time"
)

func sampleNormEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Dependency: "pkg-a", Version: "1.0.0", Date: now, Body: "  hello   world  "},
		{Dependency: "pkg-b", Version: "2.0.0", Date: now, Body: "<b>bold</b> text\n  trailing  "},
		{Dependency: "pkg-c", Version: "3.0.0", Date: now, Body: strings.Repeat("x", 200)},
	}
}

func TestNewNormalizer_NotNil(t *testing.T) {
	nr := NewNormalizer()
	if nr == nil {
		t.Fatal("expected non-nil Normalizer")
	}
}

func TestNormalizer_CollapseSpaces(t *testing.T) {
	nr := NewNormalizer()
	entries := []Entry{{Body: "  hello   world  "}}
	out := nr.Normalize(entries)
	if out[0].Body != "hello world" {
		t.Errorf("expected 'hello world', got %q", out[0].Body)
	}
}

func TestNormalizer_StripHTML(t *testing.T) {
	nr := NewNormalizer(WithStripHTML(true))
	entries := []Entry{{Body: "<b>bold</b> and <i>italic</i>"}}
	out := nr.Normalize(entries)
	if strings.Contains(out[0].Body, "<") {
		t.Errorf("expected HTML stripped, got %q", out[0].Body)
	}
	if !strings.Contains(out[0].Body, "bold") {
		t.Errorf("expected text preserved, got %q", out[0].Body)
	}
}

func TestNormalizer_NoStripHTML_ByDefault(t *testing.T) {
	nr := NewNormalizer()
	entries := []Entry{{Body: "<b>bold</b>"}}
	out := nr.Normalize(entries)
	if !strings.Contains(out[0].Body, "<b>") {
		t.Errorf("expected HTML preserved by default, got %q", out[0].Body)
	}
}

func TestNormalizer_MaxLength(t *testing.T) {
	nr := NewNormalizer(WithMaxLength(10))
	entries := []Entry{{Body: strings.Repeat("a", 50)}}
	out := nr.Normalize(entries)
	if len(out[0].Body) != 10 {
		t.Errorf("expected length 10, got %d", len(out[0].Body))
	}
}

func TestNormalizer_MaxLength_Zero_NoTruncation(t *testing.T) {
	nr := NewNormalizer(WithMaxLength(0))
	body := strings.Repeat("z", 300)
	entries := []Entry{{Body: body}}
	out := nr.Normalize(entries)
	if len(out[0].Body) != 300 {
		t.Errorf("expected no truncation, got length %d", len(out[0].Body))
	}
}

func TestNormalizer_PreservesNewlines(t *testing.T) {
	nr := NewNormalizer()
	entries := []Entry{{Body: "line one\nline two\n"}}
	out := nr.Normalize(entries)
	if !strings.Contains(out[0].Body, "\n") {
		t.Errorf("expected newlines preserved, got %q", out[0].Body)
	}
}

func TestNormalizer_EmptyBody(t *testing.T) {
	nr := NewNormalizer(WithStripHTML(true), WithMaxLength(100))
	entries := []Entry{{Body: ""}}
	out := nr.Normalize(entries)
	if out[0].Body != "" {
		t.Errorf("expected empty body, got %q", out[0].Body)
	}
}
