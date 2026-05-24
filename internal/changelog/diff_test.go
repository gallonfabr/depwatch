package changelog_test

import (
	"testing"
	"time"

	"github.com/your-org/depwatch/internal/changelog"
)

var (
	now   = time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	older = now.Add(-48 * time.Hour)
)

func sampleDiffEntries() []changelog.Entry {
	return []changelog.Entry{
		{Dependency: "libA", Version: "1.0.0", Date: older},
		{Dependency: "libA", Version: "1.1.0", Date: now},
		{Dependency: "libB", Version: "2.0.0", Date: now},
	}
}

func TestNewDiff_NotNil(t *testing.T) {
	if changelog.NewDiff() == nil {
		t.Fatal("expected non-nil Diff")
	}
}

func TestDiff_Apply_AllNew(t *testing.T) {
	d := changelog.NewDiff()
	novel := d.Apply(nil, sampleDiffEntries())
	if len(novel) != 3 {
		t.Fatalf("expected 3 novel entries, got %d", len(novel))
	}
}

func TestDiff_Apply_NoneNew(t *testing.T) {
	d := changelog.NewDiff()
	entries := sampleDiffEntries()
	novel := d.Apply(entries, entries)
	if len(novel) != 0 {
		t.Fatalf("expected 0 novel entries, got %d", len(novel))
	}
}

func TestDiff_Apply_PartiallyNew(t *testing.T) {
	d := changelog.NewDiff()
	previous := sampleDiffEntries()[:2]
	current := sampleDiffEntries()
	novel := d.Apply(previous, current)
	if len(novel) != 1 {
		t.Fatalf("expected 1 novel entry, got %d", len(novel))
	}
	if novel[0].Dependency != "libB" {
		t.Errorf("expected libB, got %s", novel[0].Dependency)
	}
}

func TestDiff_Apply_SameVersionDifferentDep(t *testing.T) {
	d := changelog.NewDiff()
	previous := []changelog.Entry{{Dependency: "libA", Version: "1.0.0"}}
	current := []changelog.Entry{{Dependency: "libB", Version: "1.0.0"}}
	novel := d.Apply(previous, current)
	if len(novel) != 1 {
		t.Fatalf("expected 1 novel entry, got %d", len(novel))
	}
}

func TestDiff_Summarise_Empty(t *testing.T) {
	d := changelog.NewDiff()
	s := d.Summarise(nil)
	if s.Added != 0 {
		t.Errorf("expected 0 added, got %d", s.Added)
	}
}

func TestDiff_Summarise_Counts(t *testing.T) {
	d := changelog.NewDiff()
	s := d.Summarise(sampleDiffEntries())
	if s.Added != 3 {
		t.Errorf("expected 3 added, got %d", s.Added)
	}
	if !s.Oldest.Equal(older) {
		t.Errorf("expected oldest %v, got %v", older, s.Oldest)
	}
	if !s.Newest.Equal(now) {
		t.Errorf("expected newest %v, got %v", now, s.Newest)
	}
}
