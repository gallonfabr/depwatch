package changelog

import (
	"testing"
)

var sampleBadgeEntries = []Entry{
	{Dependency: "react", Tags: []string{"security", "breaking"}},
	{Dependency: "lodash", Tags: []string{"feature"}},
	{Dependency: "axios", Tags: []string{}},
}

func TestNewBadger_NotNil(t *testing.T) {
	b := NewBadger()
	if b == nil {
		t.Fatal("expected non-nil Badger")
	}
}

func TestBadger_NoRules_NoBadges(t *testing.T) {
	b := NewBadger()
	out := b.Apply(sampleBadgeEntries)
	for _, e := range out {
		if len(e.Badges) != 0 {
			t.Errorf("expected no badges, got %v for %s", e.Badges, e.Dependency)
		}
	}
}

func TestBadger_SecurityTag_AttachesBadge(t *testing.T) {
	b := NewBadger(
		WithBadgeRule("security", "Security", "#e11d48"),
	)
	out := b.Apply(sampleBadgeEntries)
	if len(out[0].Badges) != 1 {
		t.Fatalf("expected 1 badge on react, got %d", len(out[0].Badges))
	}
	if out[0].Badges[0].Label != "Security" {
		t.Errorf("expected label Security, got %s", out[0].Badges[0].Label)
	}
	if out[0].Badges[0].Color != "#e11d48" {
		t.Errorf("expected color #e11d48, got %s", out[0].Badges[0].Color)
	}
}

func TestBadger_MultipleRules_MultipleBadges(t *testing.T) {
	b := NewBadger(
		WithBadgeRule("security", "Security", "#e11d48"),
		WithBadgeRule("breaking", "Breaking", "#f97316"),
	)
	out := b.Apply(sampleBadgeEntries)
	if len(out[0].Badges) != 2 {
		t.Fatalf("expected 2 badges on react, got %d", len(out[0].Badges))
	}
}

func TestBadger_NoBadge_WhenTagAbsent(t *testing.T) {
	b := NewBadger(
		WithBadgeRule("security", "Security", "#e11d48"),
	)
	out := b.Apply(sampleBadgeEntries)
	if len(out[1].Badges) != 0 {
		t.Errorf("expected no badge on lodash, got %v", out[1].Badges)
	}
	if len(out[2].Badges) != 0 {
		t.Errorf("expected no badge on axios, got %v", out[2].Badges)
	}
}

func TestBadger_NoDuplicateBadges(t *testing.T) {
	b := NewBadger(
		WithBadgeRule("security", "Security", "#e11d48"),
		WithBadgeRule("security", "Security", "#e11d48"),
	)
	out := b.Apply(sampleBadgeEntries)
	if len(out[0].Badges) != 1 {
		t.Errorf("expected exactly 1 badge, got %d", len(out[0].Badges))
	}
}

func TestBadger_DoesNotMutateInput(t *testing.T) {
	original := []Entry{
		{Dependency: "vue", Tags: []string{"feature"}},
	}
	b := NewBadger(WithBadgeRule("feature", "Feature", "#22c55e"))
	_ = b.Apply(original)
	if len(original[0].Badges) != 0 {
		t.Error("Apply must not mutate the input slice")
	}
}
