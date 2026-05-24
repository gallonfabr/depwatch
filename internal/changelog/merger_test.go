package changelog

import (
	"testing"
	"time"
)

var mergerBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func sampleMergerEntries() ([]Entry, []Entry) {
	a := []Entry{
		{Dependency: "libA", Version: "1.0.0", Date: mergerBase},
		{Dependency: "libA", Version: "1.1.0", Date: mergerBase.Add(24 * time.Hour)},
	}
	b := []Entry{
		{Dependency: "libA", Version: "1.1.0", Date: mergerBase.Add(24 * time.Hour)}, // duplicate
		{Dependency: "libB", Version: "2.0.0", Date: mergerBase},
	}
	return a, b
}

func TestNewMerger_NotNil(t *testing.T) {
	if NewMerger() == nil {
		t.Fatal("expected non-nil Merger")
	}
}

func TestMerger_Merge_CombinesSlices(t *testing.T) {
	a, b := sampleMergerEntries()
	m := NewMerger()
	result := m.Merge(a, b)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestMerger_Merge_RemovesDuplicates(t *testing.T) {
	a, b := sampleMergerEntries()
	m := NewMerger()
	result := m.Merge(a, b)
	seen := make(map[string]bool)
	for _, e := range result {
		key := e.Dependency + "@" + e.Version
		if seen[key] {
			t.Errorf("duplicate entry found: %s", key)
		}
		seen[key] = true
	}
}

func TestMerger_Merge_PreservesOrder(t *testing.T) {
	a, b := sampleMergerEntries()
	m := NewMerger()
	result := m.Merge(a, b)
	if result[0].Version != "1.0.0" {
		t.Errorf("expected first entry version 1.0.0, got %s", result[0].Version)
	}
}

func TestMerger_Merge_EmptySources(t *testing.T) {
	m := NewMerger()
	result := m.Merge()
	if result != nil && len(result) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(result))
	}
}

func TestMerger_Merge_SingleSource(t *testing.T) {
	a, _ := sampleMergerEntries()
	m := NewMerger()
	result := m.Merge(a)
	if len(result) != len(a) {
		t.Fatalf("expected %d entries, got %d", len(a), len(result))
	}
}

func TestMergeAll_ConvenienceWrapper(t *testing.T) {
	a, b := sampleMergerEntries()
	result := MergeAll(a, b)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}
