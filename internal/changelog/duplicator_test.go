package changelog

import (
	"testing"
	"time"
)

var sampleCrossEntries = []Entry{
	{Dependency: "react", Version: "18.0.0", Date: time.Now()},
	{Dependency: "react", Version: "18.0.0", Date: time.Now()}, // duplicate
	{Dependency: "react", Version: "18.1.0", Date: time.Now()},
	{Dependency: "vue", Version: "3.0.0", Date: time.Now()},
	{Dependency: "Vue", Version: "3.0.0", Date: time.Now()}, // case-insensitive duplicate
}

func TestNewCrossDeduplicator_NotNil(t *testing.T) {
	cd := NewCrossDeduplicator()
	if cd == nil {
		t.Fatal("expected non-nil CrossDeduplicator")
	}
}

func TestCrossDeduplicator_Apply_RemovesDuplicates(t *testing.T) {
	cd := NewCrossDeduplicator()
	out := cd.Apply(sampleCrossEntries)
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestCrossDeduplicator_Apply_CaseInsensitive(t *testing.T) {
	cd := NewCrossDeduplicator()
	entries := []Entry{
		{Dependency: "Django", Version: "4.0"},
		{Dependency: "django", Version: "4.0"},
	}
	out := cd.Apply(entries)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry after case-insensitive dedup, got %d", len(out))
	}
}

func TestCrossDeduplicator_Apply_CrossCallDedup(t *testing.T) {
	cd := NewCrossDeduplicator()
	first := []Entry{{Dependency: "lodash", Version: "4.17.21"}}
	second := []Entry{{Dependency: "lodash", Version: "4.17.21"}}
	cd.Apply(first)
	out := cd.Apply(second)
	if len(out) != 0 {
		t.Fatalf("expected 0 entries on second call, got %d", len(out))
	}
}

func TestCrossDeduplicator_Reset_ClearsSeen(t *testing.T) {
	cd := NewCrossDeduplicator()
	entries := []Entry{{Dependency: "axios", Version: "1.0.0"}}
	cd.Apply(entries)
	cd.Reset()
	out := cd.Apply(entries)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry after reset, got %d", len(out))
	}
}

func TestCrossDeduplicator_Apply_Empty(t *testing.T) {
	cd := NewCrossDeduplicator()
	out := cd.Apply([]Entry{})
	if len(out) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(out))
	}
}

func TestCrossDeduplicator_Apply_DifferentVersions_BothKept(t *testing.T) {
	cd := NewCrossDeduplicator()
	entries := []Entry{
		{Dependency: "express", Version: "4.18.0"},
		{Dependency: "express", Version: "4.19.0"},
	}
	out := cd.Apply(entries)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries for different versions, got %d", len(out))
	}
}
