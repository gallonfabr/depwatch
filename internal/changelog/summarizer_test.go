package changelog

import (
	"strings"
	"testing"
)

func TestNewSummarizer_NotNil(t *testing.T) {
	s := NewSummarizer()
	if s == nil {
		t.Fatal("expected non-nil Summarizer")
	}
}

func TestNewSummarizer_DefaultLength(t *testing.T) {
	s := NewSummarizer()
	if s.maxRunes != defaultSummaryLength {
		t.Fatalf("expected default %d, got %d", defaultSummaryLength, s.maxRunes)
	}
}

func TestNewSummarizer_CustomLength(t *testing.T) {
	s := NewSummarizer(WithSummaryLength(50))
	if s.maxRunes != 50 {
		t.Fatalf("expected 50, got %d", s.maxRunes)
	}
}

func TestNewSummarizer_ZeroLengthIgnored(t *testing.T) {
	s := NewSummarizer(WithSummaryLength(0))
	if s.maxRunes != defaultSummaryLength {
		t.Fatalf("zero should be ignored, got %d", s.maxRunes)
	}
}

func TestSummarizer_Apply_ShortBody_Unchanged(t *testing.T) {
	s := NewSummarizer()
	entries := []Entry{{Body: "short body"}}
	out := s.Apply(entries)
	if out[0].Body != "short body" {
		t.Fatalf("unexpected body: %q", out[0].Body)
	}
}

func TestSummarizer_Apply_CollapsesNewlines(t *testing.T) {
	s := NewSummarizer()
	entries := []Entry{{Body: "line one\nline two\nline three"}}
	out := s.Apply(entries)
	if strings.Contains(out[0].Body, "\n") {
		t.Fatalf("expected no newlines, got %q", out[0].Body)
	}
	if out[0].Body != "line one line two line three" {
		t.Fatalf("unexpected body: %q", out[0].Body)
	}
}

func TestSummarizer_Apply_TruncatesLongBody(t *testing.T) {
	s := NewSummarizer(WithSummaryLength(10))
	entries := []Entry{{Body: "hello world this is a long sentence"}}
	out := s.Apply(entries)
	runes := []rune(out[0].Body)
	// last rune should be ellipsis
	if runes[len(runes)-1] != '…' {
		t.Fatalf("expected ellipsis suffix, got %q", out[0].Body)
	}
	// visible content should be exactly maxRunes + ellipsis
	if len(runes) != 11 {
		t.Fatalf("expected 11 runes (10 + ellipsis), got %d", len(runes))
	}
}

func TestSummarizer_Apply_DoesNotMutateInput(t *testing.T) {
	s := NewSummarizer(WithSummaryLength(5))
	orig := "hello world"
	entries := []Entry{{Body: orig}}
	s.Apply(entries)
	if entries[0].Body != orig {
		t.Fatal("Apply must not mutate input slice")
	}
}

func TestSummarizer_Apply_EmptyBody(t *testing.T) {
	s := NewSummarizer()
	out := s.Apply([]Entry{{Body: ""}})
	if out[0].Body != "" {
		t.Fatalf("expected empty body, got %q", out[0].Body)
	}
}
