package changelog

import (
	"testing"
	"time"
)

func TestNewChanger_NilStore_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil store")
		}
	}()
	NewChanger(nil)
}

func TestNewChanger_Valid(t *testing.T) {
	s := NewSnapshotStore()
	c := NewChanger(s)
	if c == nil {
		t.Fatal("expected non-nil Changer")
	}
}

func TestChanger_HasChanged_NoSnapshot_TrueWhenEntries(t *testing.T) {
	s := NewSnapshotStore()
	c := NewChanger(s)
	entries := []Entry{{Dependency: "dep", Version: "v1.0.0"}}
	if !c.HasChanged("dep", entries) {
		t.Error("expected HasChanged=true when no prior snapshot")
	}
}

func TestChanger_HasChanged_NoSnapshot_FalseWhenEmpty(t *testing.T) {
	s := NewSnapshotStore()
	c := NewChanger(s)
	if c.HasChanged("dep", []Entry{}) {
		t.Error("expected HasChanged=false for empty entries with no snapshot")
	}
}

func TestChanger_HasChanged_SameEntries_False(t *testing.T) {
	s := NewSnapshotStore()
	c := NewChanger(s)
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := []Entry{{Dependency: "dep", Version: "v1.0.0", Date: now}}
	c.Commit("dep", entries)
	if c.HasChanged("dep", entries) {
		t.Error("expected HasChanged=false for identical entries")
	}
}

func TestChanger_HasChanged_DifferentVersion_True(t *testing.T) {
	s := NewSnapshotStore()
	c := NewChanger(s)
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	old := []Entry{{Dependency: "dep", Version: "v1.0.0", Date: now}}
	c.Commit("dep", old)
	newEntries := []Entry{{Dependency: "dep", Version: "v1.1.0", Date: now}}
	if !c.HasChanged("dep", newEntries) {
		t.Error("expected HasChanged=true when version differs")
	}
}

func TestChanger_HasChanged_LengthDiffers_True(t *testing.T) {
	s := NewSnapshotStore()
	c := NewChanger(s)
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	old := []Entry{{Dependency: "dep", Version: "v1.0.0", Date: now}}
	c.Commit("dep", old)
	newEntries := []Entry{
		{Dependency: "dep", Version: "v1.0.0", Date: now},
		{Dependency: "dep", Version: "v1.1.0", Date: now},
	}
	if !c.HasChanged("dep", newEntries) {
		t.Error("expected HasChanged=true when entry count differs")
	}
}

func TestChanger_Commit_UpdatesSnapshot(t *testing.T) {
	s := NewSnapshotStore()
	c := NewChanger(s)
	now := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	v1 := []Entry{{Dependency: "lib", Version: "v1.0.0", Date: now}}
	c.Commit("lib", v1)
	v2 := []Entry{{Dependency: "lib", Version: "v2.0.0", Date: now}}
	c.Commit("lib", v2)
	got, ok := s.Get("lib")
	if !ok {
		t.Fatal("expected snapshot to exist after Commit")
	}
	if got[0].Version != "v2.0.0" {
		t.Errorf("expected v2.0.0, got %s", got[0].Version)
	}
}

func TestChanger_HasChanged_MultipleKeys_Isolated(t *testing.T) {
	// Verify that committing entries for one key does not affect HasChanged for another key.
	s := NewSnapshotStore()
	c := NewChanger(s)
	now := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	entriesA := []Entry{{Dependency: "libA", Version: "v1.0.0", Date: now}}
	c.Commit("libA", entriesA)

	entriesB := []Entry{{Dependency: "libB", Version: "v2.0.0", Date: now}}
	if !c.HasChanged("libB", entriesB) {
		t.Error("expected HasChanged=true for libB which has no prior snapshot")
	}
}
