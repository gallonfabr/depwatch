package changelog

import (
	"testing"
	"time"
)

func sampleTruncatorEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Dependency: "react", Version: "18.0.0", Date: now},
		{Dependency: "react", Version: "17.0.2", Date: now},
		{Dependency: "react", Version: "17.0.1", Date: now},
		{Dependency: "lodash", Version: "4.17.21", Date: now},
		{Dependency: "lodash", Version: "4.17.20", Date: now},
		{Dependency: "axios", Version: "1.4.0", Date: now},
	}
}

func TestNewTruncator_NotNil(t *testing.T) {
	tr := NewTruncator(3)
	if tr == nil {
		t.Fatal("expected non-nil Truncator")
	}
}

func TestTruncator_ZeroMax_ReturnsAll(t *testing.T) {
	entries := sampleTruncatorEntries()
	tr := NewTruncator(0)
	got := tr.Apply(entries)
	if len(got) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(got))
	}
}

func TestTruncator_NegativeMax_ReturnsAll(t *testing.T) {
	entries := sampleTruncatorEntries()
	tr := NewTruncator(-1)
	got := tr.Apply(entries)
	if len(got) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(got))
	}
}

func TestTruncator_CapPerDependency(t *testing.T) {
	entries := sampleTruncatorEntries()
	tr := NewTruncator(2)
	got := tr.Apply(entries)

	counts := make(map[string]int)
	for _, e := range got {
		counts[e.Dependency]++
	}

	for dep, count := range counts {
		if count > 2 {
			t.Errorf("dependency %q has %d entries, want <= 2", dep, count)
		}
	}
}

func TestTruncator_PreservesOrder(t *testing.T) {
	entries := sampleTruncatorEntries()
	tr := NewTruncator(2)
	got := tr.Apply(entries)

	// First react entry kept should be 18.0.0 (appears first in input).
	for _, e := range got {
		if e.Dependency == "react" {
			if e.Version != "18.0.0" {
				t.Errorf("expected first react entry to be 18.0.0, got %s", e.Version)
			}
			break
		}
	}
}

func TestTruncator_MaxOne_KeepsFirstPerDep(t *testing.T) {
	entries := sampleTruncatorEntries()
	tr := NewTruncator(1)
	got := tr.Apply(entries)

	seen := make(map[string]bool)
	for _, e := range got {
		if seen[e.Dependency] {
			t.Errorf("duplicate dependency %q after max=1 truncation", e.Dependency)
		}
		seen[e.Dependency] = true
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 entries (one per dep), got %d", len(got))
	}
}

func TestTruncator_EmptyInput_ReturnsEmpty(t *testing.T) {
	tr := NewTruncator(5)
	got := tr.Apply([]Entry{})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(got))
	}
}
