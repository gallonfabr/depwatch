package changelog_test

import (
	"testing"
	"time"

	"github.com/depwatch/internal/changelog"
)

var sampleMatcherEntries = []changelog.Entry{
	{Version: "v1.0.0", Body: "Fix critical security vulnerability", Date: time.Now()},
	{Version: "v1.1.0", Body: "Add new feature for dashboard", Date: time.Now()},
	{Version: "v2.0.0", Body: "Breaking change: removed deprecated API", Date: time.Now()},
	{Version: "v2.1.0", Body: "Minor bug fix", Date: time.Now()},
}

func TestNewMatcher_NotNil(t *testing.T) {
	m := changelog.NewMatcher()
	if m == nil {
		t.Fatal("expected non-nil Matcher")
	}
}

func TestMatcher_NoRules_ReturnsAll(t *testing.T) {
	m := changelog.NewMatcher()
	got := m.Apply(sampleMatcherEntries)
	if len(got) != len(sampleMatcherEntries) {
		t.Fatalf("expected %d entries, got %d", len(sampleMatcherEntries), len(got))
	}
}

func TestMatcher_BodyPattern_Matches(t *testing.T) {
	m := changelog.NewMatcher(
		changelog.WithMatchRule("body", "(?i)security"),
	)
	got := m.Apply(sampleMatcherEntries)
	if len(got) != 1 {
		t.Fatalf("expected 1 match, got %d", len(got))
	}
	if got[0].Version != "v1.0.0" {
		t.Errorf("expected v1.0.0, got %s", got[0].Version)
	}
}

func TestMatcher_MultipleRules_UnionMatch(t *testing.T) {
	m := changelog.NewMatcher(
		changelog.WithMatchRule("body", "(?i)security"),
		changelog.WithMatchRule("body", "(?i)breaking"),
	)
	got := m.Apply(sampleMatcherEntries)
	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(got))
	}
}

func TestMatcher_VersionPattern_Matches(t *testing.T) {
	m := changelog.NewMatcher(
		changelog.WithMatchRule("version", `^v2`),
	)
	got := m.Apply(sampleMatcherEntries)
	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(got))
	}
}

func TestMatcher_NoMatch_ReturnsEmpty(t *testing.T) {
	m := changelog.NewMatcher(
		changelog.WithMatchRule("body", "(?i)nonexistent_xyz"),
	)
	got := m.Apply(sampleMatcherEntries)
	if len(got) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(got))
	}
}

func TestMatcher_EmptyInput_ReturnsEmpty(t *testing.T) {
	m := changelog.NewMatcher(
		changelog.WithMatchRule("body", "security"),
	)
	got := m.Apply([]changelog.Entry{})
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}
