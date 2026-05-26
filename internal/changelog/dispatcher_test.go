package changelog_test

import (
	"sort"
	"testing"
	"time"

	"github.com/depwatch/internal/changelog"
)

func sampleDispatchEntries() []changelog.Entry {
	return []changelog.Entry{
		{Dependency: "libA", Version: "1.0.0", Labels: []string{"security"}, Date: time.Now()},
		{Dependency: "libB", Version: "2.0.0", Labels: []string{"feature"}, Date: time.Now()},
		{Dependency: "libC", Version: "3.0.0", Labels: []string{"security", "breaking"}, Date: time.Now()},
	}
}

func buildTestRouter() *changelog.Router {
	return changelog.NewRouter(
		changelog.WithRouteRule("security-channel", "security"),
		changelog.WithRouteRule("feature-channel", "feature"),
		changelog.WithRouterFallback("general"),
	)
}

func TestNewDispatcher_NilRouter_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil router")
		}
	}()
	changelog.NewDispatcher(nil)
}

func TestNewDispatcher_Valid(t *testing.T) {
	d := changelog.NewDispatcher(buildTestRouter())
	if d == nil {
		t.Fatal("expected non-nil dispatcher")
	}
}

func TestDispatcher_Dispatch_SecurityChannel(t *testing.T) {
	d := changelog.NewDispatcher(buildTestRouter())
	d.Dispatch(sampleDispatchEntries())

	entries := d.EntriesFor("security-channel")
	if len(entries) != 2 {
		t.Fatalf("expected 2 security entries, got %d", len(entries))
	}
}

func TestDispatcher_Dispatch_FeatureChannel(t *testing.T) {
	d := changelog.NewDispatcher(buildTestRouter())
	d.Dispatch(sampleDispatchEntries())

	entries := d.EntriesFor("feature-channel")
	if len(entries) != 1 {
		t.Fatalf("expected 1 feature entry, got %d", len(entries))
	}
}

func TestDispatcher_Channels_ReturnAllActive(t *testing.T) {
	d := changelog.NewDispatcher(buildTestRouter())
	d.Dispatch(sampleDispatchEntries())

	chs := d.Channels()
	sort.Strings(chs)

	if len(chs) < 2 {
		t.Fatalf("expected at least 2 channels, got %d: %v", len(chs), chs)
	}
}

func TestDispatcher_Dispatch_ClearsPreviousResults(t *testing.T) {
	d := changelog.NewDispatcher(buildTestRouter())
	d.Dispatch(sampleDispatchEntries())
	d.Dispatch([]changelog.Entry{}) // empty second dispatch

	if len(d.Channels()) != 0 {
		t.Fatal("expected no channels after empty dispatch")
	}
}

func TestDispatcher_EntriesFor_MissingChannel_ReturnsNil(t *testing.T) {
	d := changelog.NewDispatcher(buildTestRouter())
	d.Dispatch(sampleDispatchEntries())

	if got := d.EntriesFor("nonexistent"); got != nil {
		t.Fatalf("expected nil for unknown channel, got %v", got)
	}
}
