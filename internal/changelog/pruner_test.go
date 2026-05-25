package changelog

import (
	"testing"
	"time"
)

func TestNewPruner_NotNil(t *testing.T) {
	p := NewPruner()
	if p == nil {
		t.Fatal("expected non-nil Pruner")
	}
}

func TestNewPruner_DefaultMaxAge(t *testing.T) {
	p := NewPruner()
	expected := 30 * 24 * time.Hour
	if p.maxAge != expected {
		t.Fatalf("expected default maxAge %v, got %v", expected, p.maxAge)
	}
}

func TestNewPruner_CustomMaxAge(t *testing.T) {
	p := NewPruner(WithMaxAge(7 * 24 * time.Hour))
	if p.maxAge != 7*24*time.Hour {
		t.Fatalf("unexpected maxAge: %v", p.maxAge)
	}
}

func TestNewPruner_ZeroMaxAge_IgnoredUsesDefault(t *testing.T) {
	p := NewPruner(WithMaxAge(0))
	if p.maxAge != 30*24*time.Hour {
		t.Fatalf("zero maxAge should be ignored, got %v", p.maxAge)
	}
}

func TestPruner_Apply_EmptyInput(t *testing.T) {
	p := NewPruner()
	result := p.Apply([]Entry{})
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(result))
	}
}

func TestPruner_Apply_KeepsRecentEntries(t *testing.T) {
	p := NewPruner(WithMaxAge(7 * 24 * time.Hour))
	entries := []Entry{
		{Dependency: "lib-a", Version: "1.0.0", Date: time.Now().UTC().Add(-24 * time.Hour)},
		{Dependency: "lib-b", Version: "2.0.0", Date: time.Now().UTC().Add(-3 * 24 * time.Hour)},
	}
	result := p.Apply(entries)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestPruner_Apply_RemovesOldEntries(t *testing.T) {
	p := NewPruner(WithMaxAge(7 * 24 * time.Hour))
	entries := []Entry{
		{Dependency: "lib-a", Version: "1.0.0", Date: time.Now().UTC().Add(-8 * 24 * time.Hour)},
		{Dependency: "lib-b", Version: "2.0.0", Date: time.Now().UTC().Add(-1 * 24 * time.Hour)},
	}
	result := p.Apply(entries)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Dependency != "lib-b" {
		t.Fatalf("expected lib-b, got %s", result[0].Dependency)
	}
}

func TestPruner_Apply_ZeroDate_AlwaysKept(t *testing.T) {
	p := NewPruner(WithMaxAge(1 * time.Hour))
	entries := []Entry{
		{Dependency: "lib-zero", Version: "0.1.0", Date: time.Time{}},
	}
	result := p.Apply(entries)
	if len(result) != 1 {
		t.Fatalf("expected zero-date entry to be kept, got %d entries", len(result))
	}
}

func TestPruner_Apply_AllOld_ReturnsEmpty(t *testing.T) {
	p := NewPruner(WithMaxAge(24 * time.Hour))
	entries := []Entry{
		{Dependency: "lib-a", Version: "1.0.0", Date: time.Now().UTC().Add(-48 * time.Hour)},
		{Dependency: "lib-b", Version: "1.1.0", Date: time.Now().UTC().Add(-72 * time.Hour)},
	}
	result := p.Apply(entries)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(result))
	}
}
