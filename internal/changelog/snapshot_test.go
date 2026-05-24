package changelog

import (
	"testing"
	"time"
)

func sampleSnapshotEntries() []Entry {
	return []Entry{
		{Dependency: "mylib", Version: "1.0.0", Date: time.Now(), Body: "initial release"},
		{Dependency: "mylib", Version: "1.1.0", Date: time.Now(), Body: "new feature"},
	}
}

func TestNewSnapshotStore_NotNil(t *testing.T) {
	s := NewSnapshotStore()
	if s == nil {
		t.Fatal("expected non-nil SnapshotStore")
	}
}

func TestSnapshotStore_Save_And_Get(t *testing.T) {
	s := NewSnapshotStore()
	entries := sampleSnapshotEntries()
	s.Save("mylib", entries)

	snap, ok := s.Get("mylib")
	if !ok {
		t.Fatal("expected snapshot to exist")
	}
	if snap.Dependency != "mylib" {
		t.Errorf("expected dependency mylib, got %s", snap.Dependency)
	}
	if len(snap.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(snap.Entries))
	}
}

func TestSnapshotStore_Get_Missing(t *testing.T) {
	s := NewSnapshotStore()
	_, ok := s.Get("nonexistent")
	if ok {
		t.Error("expected no snapshot for unknown dependency")
	}
}

func TestSnapshotStore_Save_IsolatesMutation(t *testing.T) {
	s := NewSnapshotStore()
	entries := sampleSnapshotEntries()
	s.Save("mylib", entries)

	// mutate original slice
	entries[0].Body = "mutated"

	snap, _ := s.Get("mylib")
	if snap.Entries[0].Body == "mutated" {
		t.Error("snapshot should not reflect mutation of original slice")
	}
}

func TestSnapshotStore_All_ReturnsAll(t *testing.T) {
	s := NewSnapshotStore()
	s.Save("lib-a", sampleSnapshotEntries())
	s.Save("lib-b", sampleSnapshotEntries())

	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(all))
	}
}

func TestSnapshotStore_Clear_RemovesEntry(t *testing.T) {
	s := NewSnapshotStore()
	s.Save("mylib", sampleSnapshotEntries())
	s.Clear("mylib")

	_, ok := s.Get("mylib")
	if ok {
		t.Error("expected snapshot to be removed after Clear")
	}
}

func TestSnapshotStore_CapturedAt_Set(t *testing.T) {
	before := time.Now()
	s := NewSnapshotStore()
	s.Save("mylib", sampleSnapshotEntries())
	after := time.Now()

	snap, _ := s.Get("mylib")
	if snap.CapturedAt.Before(before) || snap.CapturedAt.After(after) {
		t.Error("CapturedAt timestamp is outside expected range")
	}
}
