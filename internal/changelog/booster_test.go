package changelog

import (
	"testing"
)

var sampleBoosterEntries = []Entry{
	{Dependency: "lodash", Version: "4.17.21"},
	{Dependency: "react", Version: "18.0.0"},
	{Dependency: "express", Version: "4.18.2"},
	{Dependency: "react", Version: "17.0.2"},
	{Dependency: "axios", Version: "1.4.0"},
}

func TestNewBooster_NotNil(t *testing.T) {
	b := NewBooster()
	if b == nil {
		t.Fatal("expected non-nil Booster")
	}
}

func TestBooster_NoPriority_OrderUnchanged(t *testing.T) {
	b := NewBooster()
	got := b.Apply(sampleBoosterEntries)
	for i, e := range got {
		if e.Dependency != sampleBoosterEntries[i].Dependency {
			t.Fatalf("index %d: got %q, want %q", i, e.Dependency, sampleBoosterEntries[i].Dependency)
		}
	}
}

func TestBooster_PriorityDep_MovedFirst(t *testing.T) {
	b := NewBooster("axios")
	got := b.Apply(sampleBoosterEntries)
	if got[0].Dependency != "axios" {
		t.Fatalf("expected axios first, got %q", got[0].Dependency)
	}
}

func TestBooster_MultiplePriority_AllPromoted(t *testing.T) {
	b := NewBooster("axios", "express")
	got := b.Apply(sampleBoosterEntries)
	for i := 0; i < 2; i++ {
		if got[i].Dependency != "axios" && got[i].Dependency != "express" {
			t.Fatalf("index %d should be a priority dep, got %q", i, got[i].Dependency)
		}
	}
}

func TestBooster_RelativeOrderPreserved(t *testing.T) {
	b := NewBooster("react")
	got := b.Apply(sampleBoosterEntries)
	// react 18.0.0 should come before react 17.0.2
	var versions []string
	for _, e := range got {
		if e.Dependency == "react" {
			versions = append(versions, e.Version)
		}
	}
	if len(versions) != 2 || versions[0] != "18.0.0" || versions[1] != "17.0.2" {
		t.Fatalf("unexpected react version order: %v", versions)
	}
}

func TestBooster_EmptyInput_ReturnsEmpty(t *testing.T) {
	b := NewBooster("react")
	got := b.Apply([]Entry{})
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %d entries", len(got))
	}
}

func TestBooster_DoesNotMutateInput(t *testing.T) {
	input := []Entry{
		{Dependency: "a"}, {Dependency: "b"}, {Dependency: "c"},
	}
	orig := make([]Entry, len(input))
	copy(orig, input)
	b := NewBooster("c")
	_ = b.Apply(input)
	for i, e := range input {
		if e.Dependency != orig[i].Dependency {
			t.Fatalf("input mutated at index %d", i)
		}
	}
}
