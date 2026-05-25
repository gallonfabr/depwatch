package changelog

import (
	"testing"
)

var sampleSemVerEntries = []Entry{
	{Version: "v1.0.0", Dependency: "lib"},
	{Version: "v1.2.3", Dependency: "lib"},
	{Version: "v2.0.0", Dependency: "lib"},
	{Version: "v2.5.1", Dependency: "lib"},
	{Version: "v3.0.0", Dependency: "lib"},
	{Version: "not-a-version", Dependency: "lib"},
}

func TestNewSemVerFilter_NotNil(t *testing.T) {
	f := NewSemVerFilter()
	if f == nil {
		t.Fatal("expected non-nil SemVerFilter")
	}
}

func TestSemVerFilter_NoOptions_ReturnsAll(t *testing.T) {
	f := NewSemVerFilter()
	out := f.Apply(sampleSemVerEntries)
	if len(out) != len(sampleSemVerEntries) {
		t.Fatalf("expected %d entries, got %d", len(sampleSemVerEntries), len(out))
	}
}

func TestSemVerFilter_MinOnly(t *testing.T) {
	f := NewSemVerFilter(WithMinVersion("v2.0.0"))
	out := f.Apply(sampleSemVerEntries)
	// v2.0.0, v2.5.1, v3.0.0, not-a-version
	if len(out) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(out))
	}
}

func TestSemVerFilter_MaxOnly(t *testing.T) {
	f := NewSemVerFilter(WithMaxVersion("v2.0.0"))
	out := f.Apply(sampleSemVerEntries)
	// v1.0.0, v1.2.3, v2.0.0, not-a-version
	if len(out) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(out))
	}
}

func TestSemVerFilter_MinAndMax(t *testing.T) {
	f := NewSemVerFilter(WithMinVersion("v1.2.3"), WithMaxVersion("v2.5.1"))
	out := f.Apply(sampleSemVerEntries)
	// v1.2.3, v2.0.0, v2.5.1, not-a-version
	if len(out) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(out))
	}
}

func TestSemVerFilter_UnparseableVersionPassesThrough(t *testing.T) {
	f := NewSemVerFilter(WithMinVersion("v9.0.0"))
	out := f.Apply(sampleSemVerEntries)
	// only not-a-version passes (no parseable version >= 9.0.0)
	if len(out) != 1 || out[0].Version != "not-a-version" {
		t.Fatalf("expected only unparseable entry, got %+v", out)
	}
}

func TestSemVerFilter_String_NoBounds(t *testing.T) {
	f := NewSemVerFilter()
	s := f.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}

func TestSemVerFilter_String_WithBounds(t *testing.T) {
	f := NewSemVerFilter(WithMinVersion("v1.0.0"), WithMaxVersion("v2.0.0"))
	s := f.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}

func TestSemVerFilter_InvalidMinIgnored(t *testing.T) {
	f := NewSemVerFilter(WithMinVersion("not-valid"))
	if f.min != nil {
		t.Fatal("expected min to be nil for invalid version")
	}
}
