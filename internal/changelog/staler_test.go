package changelog

import (
	"testing"
	"time"
)

func TestNewStaler_NotNil(t *testing.T) {
	s := NewStaler(24 * time.Hour)
	if s == nil {
		t.Fatal("expected non-nil Staler")
	}
}

func TestStaler_ZeroThreshold_NoTagging(t *testing.T) {
	s := NewStaler(0)
	entries := []Entry{
		{Dependency: "lib", Date: time.Now().Add(-72 * time.Hour)},
	}
	out := s.Apply(entries)
	if len(out[0].Tags) != 0 {
		t.Errorf("expected no tags with zero threshold, got %v", out[0].Tags)
	}
}

func TestStaler_OldEntry_TaggedStale(t *testing.T) {
	s := NewStaler(24 * time.Hour)
	old := time.Now().UTC().Add(-48 * time.Hour)
	entries := []Entry{
		{Dependency: "lib", Date: old},
	}
	out := s.Apply(entries)
	if !containsTagStr(out[0].Tags, "stale") {
		t.Errorf("expected 'stale' tag on old entry, got %v", out[0].Tags)
	}
}

func TestStaler_RecentEntry_NotTagged(t *testing.T) {
	s := NewStaler(24 * time.Hour)
	recent := time.Now().UTC().Add(-1 * time.Hour)
	entries := []Entry{
		{Dependency: "lib", Date: recent},
	}
	out := s.Apply(entries)
	if containsTagStr(out[0].Tags, "stale") {
		t.Error("expected recent entry NOT to be tagged stale")
	}
}

func TestStaler_ZeroDate_NotTagged(t *testing.T) {
	s := NewStaler(24 * time.Hour)
	entries := []Entry{
		{Dependency: "lib", Date: time.Time{}},
	}
	out := s.Apply(entries)
	if containsTagStr(out[0].Tags, "stale") {
		t.Error("expected zero-date entry NOT to be tagged stale")
	}
}

func TestStaler_NoDuplicateStaleTag(t *testing.T) {
	s := NewStaler(24 * time.Hour)
	old := time.Now().UTC().Add(-48 * time.Hour)
	entries := []Entry{
		{Dependency: "lib", Date: old, Tags: []string{"stale"}},
	}
	out := s.Apply(entries)
	count := 0
	for _, tg := range out[0].Tags {
		if tg == "stale" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly 1 'stale' tag, got %d", count)
	}
}

func TestStaler_MixedEntries(t *testing.T) {
	s := NewStaler(24 * time.Hour)
	entries := []Entry{
		{Dependency: "old-lib", Date: time.Now().UTC().Add(-48 * time.Hour)},
		{Dependency: "new-lib", Date: time.Now().UTC().Add(-1 * time.Hour)},
	}
	out := s.Apply(entries)
	if !containsTagStr(out[0].Tags, "stale") {
		t.Error("expected first entry to be stale")
	}
	if containsTagStr(out[1].Tags, "stale") {
		t.Error("expected second entry NOT to be stale")
	}
}

// containsTagStr is a helper used only in staler tests.
func containsTagStr(tags []string, target string) bool {
	for _, t := range tags {
		if t == target {
			return true
		}
	}
	return false
}
