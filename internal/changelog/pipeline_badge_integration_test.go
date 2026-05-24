package changelog_test

import (
	"testing"

	"github.com/yourorg/depwatch/internal/changelog"
)

func TestPipeline_WithBadger_AttachesBadges(t *testing.T) {
	entries := []changelog.Entry{
		{Dependency: "react", Version: "18.0.0", Tags: []string{"security"}},
		{Dependency: "lodash", Version: "4.17.21", Tags: []string{"feature"}},
	}

	badger := changelog.NewBadger(
		changelog.WithBadgeRule("security", "Security", "#e11d48"),
	)

	pipeline := changelog.NewPipeline(
		changelog.NewDeduplicator(),
		badger,
	)

	out := pipeline.Run(entries)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}

	var react changelog.Entry
	for _, e := range out {
		if e.Dependency == "react" {
			react = e
		}
	}
	if len(react.Badges) != 1 {
		t.Fatalf("expected 1 badge on react, got %d", len(react.Badges))
	}
	if react.Badges[0].Label != "Security" {
		t.Errorf("expected Security badge, got %s", react.Badges[0].Label)
	}
}

func TestPipeline_WithBadger_NoMatchNoMutation(t *testing.T) {
	entries := []changelog.Entry{
		{Dependency: "axios", Version: "1.0.0", Tags: []string{}},
	}

	badger := changelog.NewBadger(
		changelog.WithBadgeRule("security", "Security", "#e11d48"),
	)

	pipeline := changelog.NewPipeline(badger)
	out := pipeline.Run(entries)

	if len(out[0].Badges) != 0 {
		t.Errorf("expected no badges, got %v", out[0].Badges)
	}
}

func TestPipeline_WithBadger_EmptyEntries_DoesNotPanic(t *testing.T) {
	badger := changelog.NewBadger(
		changelog.WithBadgeRule("breaking", "Breaking", "#f97316"),
	)
	pipeline := changelog.NewPipeline(badger)
	out := pipeline.Run([]changelog.Entry{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
