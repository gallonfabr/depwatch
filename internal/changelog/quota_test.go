package changelog

import (
	"testing"
	"time"
)

func makeQuotaEntries(dep string, n int) []Entry {
	out := make([]Entry, n)
	for i := range out {
		out[i] = Entry{Dependency: dep, Version: "1.0.0"}
	}
	return out
}

func TestNewQuota_NotNil(t *testing.T) {
	q := NewQuota(5, time.Minute)
	if q == nil {
		t.Fatal("expected non-nil Quota")
	}
}

func TestNewQuota_PanicsOnZeroMax(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for max < 1")
		}
	}()
	NewQuota(0, time.Minute)
}

func TestNewQuota_PanicsOnNegativeWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for non-positive window")
		}
	}()
	NewQuota(1, -time.Second)
}

func TestQuota_Apply_UnderLimit_AllAllowed(t *testing.T) {
	q := NewQuota(5, time.Minute)
	entries := makeQuotaEntries("dep-a", 3)
	got := q.Apply(entries)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestQuota_Apply_ExceedsLimit_CapsOutput(t *testing.T) {
	q := NewQuota(2, time.Minute)
	entries := makeQuotaEntries("dep-b", 5)
	got := q.Apply(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries (cap), got %d", len(got))
	}
}

func TestQuota_Apply_IndependentPerDependency(t *testing.T) {
	q := NewQuota(2, time.Minute)
	a := makeQuotaEntries("dep-a", 2)
	b := makeQuotaEntries("dep-b", 2)
	got := q.Apply(append(a, b...))
	if len(got) != 4 {
		t.Fatalf("expected 4 entries across two deps, got %d", len(got))
	}
}

func TestQuota_Reset_ClearsState(t *testing.T) {
	q := NewQuota(2, time.Minute)
	_ = q.Apply(makeQuotaEntries("dep-c", 2))
	q.Reset()
	got := q.Apply(makeQuotaEntries("dep-c", 2))
	if len(got) != 2 {
		t.Fatalf("expected 2 entries after reset, got %d", len(got))
	}
}

func TestQuota_Apply_EmptyInput_ReturnsEmpty(t *testing.T) {
	q := NewQuota(5, time.Minute)
	got := q.Apply([]Entry{})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}
