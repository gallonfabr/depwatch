package changelog_test

import (
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/changelog"
)

func TestPipeline_Pruner_RemovesStaleEntries(t *testing.T) {
	now := time.Now().UTC()
	entries := []changelog.Entry{
		{Dependency: "alpha", Version: "1.0.0", Date: now.Add(-2 * 24 * time.Hour)},
		{Dependency: "beta", Version: "2.0.0", Date: now.Add(-10 * 24 * time.Hour)},
		{Dependency: "gamma", Version: "3.0.0", Date: now.Add(-1 * 24 * time.Hour)},
	}

	pruner := changelog.NewPruner(changelog.WithMaxAge(7 * 24 * time.Hour))
	pipeline := changelog.NewPipeline(pruner)
	result := pipeline.Run(entries)

	if len(result) != 2 {
		t.Fatalf("expected 2 entries after pruning, got %d", len(result))
	}
	for _, e := range result {
		if e.Dependency == "beta" {
			t.Errorf("beta should have been pruned as too old")
		}
	}
}

func TestPipeline_Pruner_AllRecent_NoneRemoved(t *testing.T) {
	now := time.Now().UTC()
	entries := []changelog.Entry{
		{Dependency: "alpha", Version: "1.0.0", Date: now.Add(-1 * time.Hour)},
		{Dependency: "beta", Version: "2.0.0", Date: now.Add(-2 * time.Hour)},
	}

	pruner := changelog.NewPruner(changelog.WithMaxAge(24 * time.Hour))
	pipeline := changelog.NewPipeline(pruner)
	result := pipeline.Run(entries)

	if len(result) != 2 {
		t.Fatalf("expected all 2 entries retained, got %d", len(result))
	}
}

func TestPipeline_Pruner_EmptyInput_DoesNotPanic(t *testing.T) {
	pruner := changelog.NewPruner()
	pipeline := changelog.NewPipeline(pruner)
	result := pipeline.Run([]changelog.Entry{})
	if result == nil {
		t.Fatal("expected non-nil result for empty input")
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(result))
	}
}
