package changelog

import (
	"testing"
	"time"
)

var sampleSuppressorEntries = []Entry{
	{Dependency: "react", Version: "18.0.0", Date: time.Now()},
	{Dependency: "lodash", Version: "4.17.21", Date: time.Now()},
	{Dependency: "axios", Version: "1.4.0", Date: time.Now()},
	{Dependency: "React", Version: "18.1.0", Date: time.Now()}, // duplicate, different case
}

func TestNewSuppressor_NotNil(t *testing.T) {
	s := NewSuppressor(nil)
	if s == nil {
		t.Fatal("expected non-nil Suppressor")
	}
}

func TestSuppressor_NoList_ReturnsAll(t *testing.T) {
	s := NewSuppressor(nil)
	out := s.Apply(sampleSuppressorEntries)
	if len(out) != len(sampleSuppressorEntries) {
		t.Fatalf("expected %d entries, got %d", len(sampleSuppressorEntries), len(out))
	}
}

func TestSuppressor_MatchedDep_Removed(t *testing.T) {
	s := NewSuppressor([]string{"lodash"})
	out := s.Apply(sampleSuppressorEntries)
	for _, e := range out {
		if e.Dependency == "lodash" {
			t.Errorf("expected lodash to be suppressed")
		}
	}
}

func TestSuppressor_CaseInsensitiveMatch(t *testing.T) {
	s := NewSuppressor([]string{"REACT"})
	out := s.Apply(sampleSuppressorEntries)
	for _, e := range out {
		if e.Dependency == "react" || e.Dependency == "React" {
			t.Errorf("expected react (any case) to be suppressed, got %s", e.Dependency)
		}
	}
}

func TestSuppressor_MultipleEntries_AllRemoved(t *testing.T) {
	s := NewSuppressor([]string{"react", "axios"})
	out := s.Apply(sampleSuppressorEntries)
	// only lodash should remain
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Dependency != "lodash" {
		t.Errorf("expected lodash, got %s", out[0].Dependency)
	}
}

func TestSuppressor_EmptyInput_DoesNotPanic(t *testing.T) {
	s := NewSuppressor([]string{"react"})
	out := s.Apply([]Entry{})
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d entries", len(out))
	}
}

func TestSuppressor_IsSuppressed_True(t *testing.T) {
	s := NewSuppressor([]string{"lodash"})
	if !s.IsSuppressed("lodash") {
		t.Error("expected lodash to be suppressed")
	}
}

func TestSuppressor_IsSuppressed_False(t *testing.T) {
	s := NewSuppressor([]string{"lodash"})
	if s.IsSuppressed("react") {
		t.Error("expected react not to be suppressed")
	}
}
