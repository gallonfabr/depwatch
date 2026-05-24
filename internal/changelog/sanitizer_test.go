package changelog

import (
	"strings"
	"testing"
	"time"
)

var sampleSanitizeEntries = []Entry{
	{Dependency: "pkg-a", Version: "1.0.0", Date: time.Now(), Body: "Hello\r\nWorld"},
	{Dependency: "pkg-b", Version: "2.0.0", Date: time.Now(), Body: "Tab\there"},
	{Dependency: "pkg-c", Version: "3.0.0", Date: time.Now(), Body: "ctrl\x01chars\x1F here"},
	{Dependency: "pkg-d", Version: "4.0.0", Date: time.Now(), Body: strings.Repeat("a", 300)},
}

func TestNewSanitizer_NotNil(t *testing.T) {
	if NewSanitizer() == nil {
		t.Fatal("expected non-nil Sanitizer")
	}
}

func TestSanitizer_NormalisesWindowsLineEndings(t *testing.T) {
	s := NewSanitizer()
	result := s.Apply([]Entry{{Body: "line1\r\nline2\r\nline3"}})
	if strings.Contains(result[0].Body, "\r") {
		t.Errorf("expected no carriage returns, got: %q", result[0].Body)
	}
	if !strings.Contains(result[0].Body, "\n") {
		t.Errorf("expected newlines to be preserved")
	}
}

func TestSanitizer_RemovesControlChars(t *testing.T) {
	s := NewSanitizer()
	result := s.Apply([]Entry{{Body: "ctrl\x01chars\x1F here"}})
	got := result[0].Body
	if strings.ContainsAny(got, "\x01\x1F") {
		t.Errorf("control chars not removed, got: %q", got)
	}
	if !strings.Contains(got, "chars") {
		t.Errorf("expected visible text preserved, got: %q", got)
	}
}

func TestSanitizer_PreservesTabsAndNewlines(t *testing.T) {
	s := NewSanitizer()
	result := s.Apply([]Entry{{Body: "col1\tcol2\nrow2"}})
	got := result[0].Body
	if !strings.Contains(got, "\t") {
		t.Errorf("expected tab preserved, got: %q", got)
	}
	if !strings.Contains(got, "\n") {
		t.Errorf("expected newline preserved, got: %q", got)
	}
}

func TestSanitizer_TruncatesAtMaxRunes(t *testing.T) {
	s := NewSanitizer(WithMaxRunes(50))
	long := strings.Repeat("x", 300)
	result := s.Apply([]Entry{{Body: long}})
	got := result[0].Body
	if len([]rune(got)) > 50 {
		t.Errorf("expected at most 50 runes, got %d", len([]rune(got)))
	}
}

func TestSanitizer_NoTruncation_WhenUnderLimit(t *testing.T) {
	s := NewSanitizer(WithMaxRunes(200))
	body := strings.Repeat("b", 100)
	result := s.Apply([]Entry{{Body: body}})
	if len([]rune(result[0].Body)) != 100 {
		t.Errorf("expected 100 runes, got %d", len([]rune(result[0].Body)))
	}
}

func TestSanitizer_EmptyBody(t *testing.T) {
	s := NewSanitizer()
	result := s.Apply([]Entry{{Body: ""}})
	if result[0].Body != "" {
		t.Errorf("expected empty body, got %q", result[0].Body)
	}
}

func TestSanitizer_Apply_PreservesOtherFields(t *testing.T) {
	s := NewSanitizer()
	e := Entry{Dependency: "mypkg", Version: "1.2.3", Body: "ok"}
	result := s.Apply([]Entry{e})
	if result[0].Dependency != "mypkg" || result[0].Version != "1.2.3" {
		t.Errorf("non-body fields mutated: %+v", result[0])
	}
}
