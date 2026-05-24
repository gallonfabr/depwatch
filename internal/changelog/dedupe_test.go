package changelog

import (
	"testing"
	"time"
)

func sampleDedupEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Version: "v1.2.0", Date: now, Body: "feat: something"},
		{Version: "v1.1.0", Date: now.Add(-24 * time.Hour), Body: "fix: bug"},
		{Version: "v1.2.0", Date: now, Body: "feat: something"},
	}
}

func TestNewDeduplicator_NotNil(t *testing.T) {
	d := NewDeduplicator()
	if d == nil {
		t.Fatal("expected non-nil Deduplicator")
	}
}

func TestDeduplicator_Apply_RemovesDuplicates(t *testing.T) {
	d := NewDeduplicator()
	entries := sampleDedupEntries()
	result := d.Apply("mylib", entries)
	if len(result) != 2 {
		t.Fatalf("expected 2 unique entries, got %d", len(result))
	}
}

func TestDeduplicator_Apply_CrossCallDedup(t *testing.T) {
	d := NewDeduplicator()
	entries := sampleDedupEntries()[:1]

	// First call records v1.2.0
	d.Apply("mylib", entries)

	// Second call with same entry should return nothing
	result := d.Apply("mylib", entries)
	if len(result) != 0 {
		t.Fatalf("expected 0 entries on second call, got %d", len(result))
	}
}

func TestDeduplicator_Apply_DifferentDeps(t *testing.T) {
	d := NewDeduplicator()
	entries := sampleDedupEntries()[:1]

	r1 := d.Apply("libA", entries)
	r2 := d.Apply("libB", entries)

	if len(r1) != 1 || len(r2) != 1 {
		t.Fatalf("expected 1 entry each for different deps, got %d and %d", len(r1), len(r2))
	}
}

func TestDeduplicator_Apply_Empty(t *testing.T) {
	d := NewDeduplicator()
	result := d.Apply("mylib", []Entry{})
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestDeduplicator_Reset_ClearsSeen(t *testing.T) {
	d := NewDeduplicator()
	entries := sampleDedupEntries()[:1]
	d.Apply("mylib", entries)

	d.Reset()

	if d.Len() != 0 {
		t.Fatalf("expected Len 0 after Reset, got %d", d.Len())
	}

	result := d.Apply("mylib", entries)
	if len(result) != 1 {
		t.Fatalf("expected entry after reset, got %d", len(result))
	}
}

func TestDeduplicator_Len(t *testing.T) {
	d := NewDeduplicator()
	entries := sampleDedupEntries()
	d.Apply("mylib", entries)

	if d.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", d.Len())
	}
}
