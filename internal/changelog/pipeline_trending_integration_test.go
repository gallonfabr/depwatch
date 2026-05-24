package changelog_test

import (
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/changelog"
)

func TestPipeline_TrendingAfterScoring(t *testing.T) {
	now := time.Now()
	entries := []changelog.Entry{
		{Dependency: "react", Version: "18.0.0", Body: "security fix", Date: now},
		{Dependency: "react", Version: "17.0.0", Body: "bug fix", Date: now.Add(-time.Hour)},
		{Dependency: "lodash", Version: "4.17.21", Body: "security patch critical", Date: now},
		{Dependency: "lodash", Version: "4.17.20", Body: "performance improvement", Date: now.Add(-time.Hour)},
	}

	scorer := changelog.NewScorer(changelog.WithKeywords([]string{"security", "critical"}))

	// Apply scorer via pipeline
	pipe := changelog.NewPipeline(scorer)
	scored := pipe.Run(entries)

	tr := changelog.NewTrending(0)
	result := tr.Analyse(scored)

	if len(result) == 0 {
		t.Fatal("expected trending results")
	}
	// lodash has two keyword matches across two entries; react has one
	if result[0].Dependency != "lodash" {
		t.Errorf("expected lodash to trend first, got %s", result[0].Dependency)
	}
}

func TestTrending_TopN_LimitRespected(t *testing.T) {
	now := time.Now()
	entries := []changelog.Entry{
		{Dependency: "a", Score: 30, Date: now},
		{Dependency: "b", Score: 20, Date: now},
		{Dependency: "c", Score: 10, Date: now},
		{Dependency: "d", Score: 5, Date: now},
	}

	tr := changelog.NewTrending(2)
	result := tr.Analyse(entries)

	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	if result[0].Dependency != "a" || result[1].Dependency != "b" {
		t.Errorf("unexpected order: %v", result)
	}
}

func TestTrending_ZeroTopN_ReturnsAll(t *testing.T) {
	now := time.Now()
	entries := []changelog.Entry{
		{Dependency: "x", Score: 1, Date: now},
		{Dependency: "y", Score: 2, Date: now},
		{Dependency: "z", Score: 3, Date: now},
	}

	tr := changelog.NewTrending(0)
	result := tr.Analyse(entries)

	if len(result) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result))
	}
}
