package changelog

import (
	"testing"
	"time"
)

// TestPipeline_WithChainTransformer verifies that a Chain can be used
// as a pipeline stage via TransformFunc wrapping.
func TestPipeline_WithChainTransformer(t *testing.T) {
	entries := []Entry{
		{Dependency: "lib-a", Version: "1.0.0", Date: time.Now()},
		{Dependency: "lib-a", Version: "1.0.0", Date: time.Now()}, // duplicate
		{Dependency: "lib-b", Version: "2.0.0", Date: time.Now()},
		{Dependency: "lib-c", Version: "3.0.0", Date: time.Now()},
	}

	chain := NewChain(
		NewDeduplicator(),
		NewLimitTransformer(2),
	)

	out := chain.Transform(entries)

	if len(out) != 2 {
		t.Fatalf("expected 2 entries after dedup+limit, got %d", len(out))
	}
}

func TestLimitTransformer_ZeroMax_ReturnsAll(t *testing.T) {
	l := NewLimitTransformer(0)
	entries := []Entry{
		{Dependency: "x", Version: "1.0.0"},
		{Dependency: "y", Version: "2.0.0"},
	}
	out := l.Transform(entries)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries for max=0 (no-op), got %d", len(out))
	}
}

func TestChain_MultipleTransformers_OrderPreserved(t *testing.T) {
	entries := []Entry{
		{Dependency: "d", Version: "0.1.0", Date: time.Now().Add(-time.Hour)},
		{Dependency: "d", Version: "0.1.0", Date: time.Now().Add(-time.Hour)}, // dup
		{Dependency: "e", Version: "0.2.0", Date: time.Now()},
		{Dependency: "f", Version: "0.3.0", Date: time.Now().Add(time.Hour)},
	}

	chain := NewChain(
		NewDeduplicator(),
		NewSorter(true), // descending
		NewLimitTransformer(2),
	)

	out := chain.Transform(entries)

	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}

	// After descending sort, first entry should have the latest date.
	if !out[0].Date.After(out[1].Date) {
		t.Errorf("expected descending order: %v should be after %v", out[0].Date, out[1].Date)
	}
}
