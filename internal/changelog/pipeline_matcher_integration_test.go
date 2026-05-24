package changelog_test

import (
	"testing"
	"time"

	"github.com/depwatch/internal/changelog"
)

func TestPipeline_WithMatcher_FiltersEntries(t *testing.T) {
	entries := []changelog.Entry{
		{Version: "v1.0.0", Body: "Security patch applied", Date: time.Now()},
		{Version: "v1.1.0", Body: "New dashboard feature", Date: time.Now()},
		{Version: "v1.2.0", Body: "Bug fix for login", Date: time.Now()},
	}

	matcher := changelog.NewMatcher(
		changelog.WithMatchRule("body", "(?i)security"),
	)

	pipeline := changelog.NewPipeline(
		changelog.NewDeduplicator(),
		matcher,
	)

	got := pipeline.Run(entries)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry after pipeline, got %d", len(got))
	}
	if got[0].Version != "v1.0.0" {
		t.Errorf("expected v1.0.0, got %s", got[0].Version)
	}
}

func TestPipeline_WithMatcher_AllMatch_NoneDropped(t *testing.T) {
	entries := []changelog.Entry{
		{Version: "v2.0.0", Body: "Breaking: removed API", Date: time.Now()},
		{Version: "v2.1.0", Body: "Breaking: changed schema", Date: time.Now()},
	}

	matcher := changelog.NewMatcher(
		changelog.WithMatchRule("body", "(?i)breaking"),
	)

	pipeline := changelog.NewPipeline(matcher)
	got := pipeline.Run(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestPipeline_WithMatcher_EmptyResult_DoesNotPanic(t *testing.T) {
	entries := []changelog.Entry{
		{Version: "v1.0.0", Body: "routine update", Date: time.Now()},
	}

	matcher := changelog.NewMatcher(
		changelog.WithMatchRule("body", "(?i)critical"),
	)

	pipeline := changelog.NewPipeline(matcher)
	got := pipeline.Run(entries)
	if got == nil {
		t.Fatal("expected non-nil slice, got nil")
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}
