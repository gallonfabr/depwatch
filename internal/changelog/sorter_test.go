package changelog_test

import (
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/changelog"
)

var (
	now   = time.Now()
	old   = now.Add(-48 * time.Hour)
	older = now.Add(-96 * time.Hour)
)

func sampleSortEntries() []changelog.Entry {
	return []changelog.Entry{
		{Version: "v1.0.0", Date: old, Body: "old"},
		{Version: "v1.2.0", Date: now, Body: "new"},
		{Version: "v0.9.0", Date: older, Body: "oldest"},
	}
}

func TestNewSorter_NotNil(t *testing.T) {
	s := changelog.NewSorter(changelog.SortDescending)
	if s == nil {
		t.Fatal("expected non-nil Sorter")
	}
}

func TestSorter_Descending_Order(t *testing.T) {
	s := changelog.NewSorter(changelog.SortDescending)
	out := s.Apply(sampleSortEntries())

	if out[0].Version != "v1.2.0" {
		t.Errorf("expected newest first, got %s", out[0].Version)
	}
	if out[2].Version != "v0.9.0" {
		t.Errorf("expected oldest last, got %s", out[2].Version)
	}
}

func TestSorter_Ascending_Order(t *testing.T) {
	s := changelog.NewSorter(changelog.SortAscending)
	out := s.Apply(sampleSortEntries())

	if out[0].Version != "v0.9.0" {
		t.Errorf("expected oldest first, got %s", out[0].Version)
	}
	if out[2].Version != "v1.2.0" {
		t.Errorf("expected newest last, got %s", out[2].Version)
	}
}

func TestSorter_ZeroDate_PlacedLast(t *testing.T) {
	entries := []changelog.Entry{
		{Version: "v1.0.0", Date: time.Time{}, Body: "no date"},
		{Version: "v1.1.0", Date: now, Body: "has date"},
	}
	s := changelog.NewSorter(changelog.SortDescending)
	out := s.Apply(entries)

	if out[0].Version != "v1.1.0" {
		t.Errorf("expected dated entry first, got %s", out[0].Version)
	}
	if out[1].Version != "v1.0.0" {
		t.Errorf("expected zero-date entry last, got %s", out[1].Version)
	}
}

func TestSorter_DoesNotMutateInput(t *testing.T) {
	input := sampleSortEntries()
	origFirst := input[0].Version

	s := changelog.NewSorter(changelog.SortDescending)
	s.Apply(input)

	if input[0].Version != origFirst {
		t.Error("Apply mutated the original input slice")
	}
}

func TestSorter_EmptyInput(t *testing.T) {
	s := changelog.NewSorter(changelog.SortDescending)
	out := s.Apply([]changelog.Entry{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
