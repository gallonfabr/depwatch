package changelog

import (
	"testing"
	"time"
)

var samplePinnedEntries = []Entry{
	{Dependency: "react", Version: "18.2.0"},
	{Dependency: "react", Version: "18.3.0"},
	{Dependency: "lodash", Version: "4.17.21"},
	{Dependency: "axios", Version: "1.6.0"},
}

func TestNewPinner_NotNil(t *testing.T) {
	p := NewPinner(nil)
	if p == nil {
		t.Fatal("expected non-nil Pinner")
	}
}

func TestPinner_NoPins_ReturnsAll(t *testing.T) {
	p := NewPinner(nil)
	got := p.Apply(samplePinnedEntries)
	if len(got) != len(samplePinnedEntries) {
		t.Fatalf("expected %d entries, got %d", len(samplePinnedEntries), len(got))
	}
}

func TestPinner_PinnedVersion_Excluded(t *testing.T) {
	pins := []PinnedEntry{
		{Dependency: "react", Version: "18.2.0", PinnedAt: time.Now()},
	}
	p := NewPinner(pins)
	got := p.Apply(samplePinnedEntries)
	for _, e := range got {
		if e.Dependency == "react" && e.Version == "18.2.0" {
			t.Error("pinned entry react@18.2.0 should have been excluded")
		}
	}
}

func TestPinner_NonPinnedVersion_Included(t *testing.T) {
	pins := []PinnedEntry{
		{Dependency: "react", Version: "18.2.0", PinnedAt: time.Now()},
	}
	p := NewPinner(pins)
	got := p.Apply(samplePinnedEntries)
	found := false
	for _, e := range got {
		if e.Dependency == "react" && e.Version == "18.3.0" {
			found = true
		}
	}
	if !found {
		t.Error("expected react@18.3.0 to be included")
	}
}

func TestPinner_MultiplePins(t *testing.T) {
	pins := []PinnedEntry{
		{Dependency: "react", Version: "18.2.0", PinnedAt: time.Now()},
		{Dependency: "lodash", Version: "4.17.21", PinnedAt: time.Now()},
	}
	p := NewPinner(pins)
	got := p.Apply(samplePinnedEntries)
	// react@18.2.0 and lodash@4.17.21 removed; react@18.3.0 and axios@1.6.0 remain
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestIsPinned_True(t *testing.T) {
	pins := []PinnedEntry{
		{Dependency: "axios", Version: "1.6.0", PinnedAt: time.Now()},
	}
	p := NewPinner(pins)
	if !p.IsPinned("axios", "1.6.0") {
		t.Error("expected axios@1.6.0 to be pinned")
	}
}

func TestIsPinned_False_WrongVersion(t *testing.T) {
	pins := []PinnedEntry{
		{Dependency: "axios", Version: "1.6.0", PinnedAt: time.Now()},
	}
	p := NewPinner(pins)
	if p.IsPinned("axios", "1.7.0") {
		t.Error("axios@1.7.0 should not be pinned")
	}
}

func TestNewPinner_EmptyDependencyIgnored(t *testing.T) {
	pins := []PinnedEntry{
		{Dependency: "", Version: "1.0.0", PinnedAt: time.Now()},
	}
	p := NewPinner(pins)
	if len(p.pins) != 0 {
		t.Error("empty dependency should not be stored")
	}
}
