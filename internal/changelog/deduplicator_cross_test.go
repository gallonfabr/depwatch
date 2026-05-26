package changelog

import (
	"testing"
)

func TestNewCrossSourceDeduplicator_NotNil(t *testing.T) {
	d := NewCrossSourceDeduplicator()
	if d == nil {
		t.Fatal("expected non-nil CrossSourceDeduplicator")
	}
}

func TestCrossSourceDeduplicator_Apply_AllUnique(t *testing.T) {
	d := NewCrossSourceDeduplicator()
	entries := []Entry{
		{Dependency: "lib-a", Version: "1.0.0"},
		{Dependency: "lib-b", Version: "1.0.0"},
		{Dependency: "lib-a", Version: "2.0.0"},
	}
	out := d.Apply(entries)
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestCrossSourceDeduplicator_Apply_RemovesDuplicate(t *testing.T) {
	d := NewCrossSourceDeduplicator()
	entries := []Entry{
		{Dependency: "lib-a", Version: "1.0.0", Body: "first"},
		{Dependency: "lib-a", Version: "1.0.0", Body: "duplicate"},
	}
	out := d.Apply(entries)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Body != "first" {
		t.Errorf("expected first entry to be kept, got body=%q", out[0].Body)
	}
}

func TestCrossSourceDeduplicator_Apply_CrossCallDedup(t *testing.T) {
	d := NewCrossSourceDeduplicator()
	first := []Entry{{Dependency: "lib-a", Version: "1.0.0"}}
	second := []Entry{{Dependency: "lib-a", Version: "1.0.0"}}

	d.Apply(first)
	out := d.Apply(second)
	if len(out) != 0 {
		t.Fatalf("expected 0 entries on second call, got %d", len(out))
	}
}

func TestCrossSourceDeduplicator_Reset_ClearsSeen(t *testing.T) {
	d := NewCrossSourceDeduplicator()
	entries := []Entry{{Dependency: "lib-a", Version: "1.0.0"}}

	d.Apply(entries)
	d.Reset()

	out := d.Apply(entries)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry after reset, got %d", len(out))
	}
}

func TestCrossSourceDeduplicator_Apply_EmptyInput(t *testing.T) {
	d := NewCrossSourceDeduplicator()
	out := d.Apply([]Entry{})
	if len(out) != 0 {
		t.Fatalf("expected 0 entries for empty input, got %d", len(out))
	}
}
