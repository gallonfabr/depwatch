package changelog

import (
	"testing"
)

var sampleRouterEntries = []Entry{
	{Dependency: "libA", Labels: []string{"security"}},
	{Dependency: "libB", Labels: []string{"feature"}},
	{Dependency: "libC", Labels: []string{"security", "breaking"}},
	{Dependency: "libD", Labels: []string{}},
}

func TestNewRouter_NotNil(t *testing.T) {
	r := NewRouter()
	if r == nil {
		t.Fatal("expected non-nil Router")
	}
}

func TestRouter_DefaultFallback(t *testing.T) {
	r := NewRouter()
	out := r.Route(sampleRouterEntries)
	if _, ok := out["default"]; !ok {
		t.Error("expected fallback channel 'default' to exist")
	}
}

func TestRouter_CustomFallback(t *testing.T) {
	r := NewRouter(WithRouterFallback("misc"))
	out := r.Route([]Entry{{Dependency: "libX", Labels: []string{}}})
	if len(out["misc"]) != 1 {
		t.Errorf("expected 1 entry in 'misc', got %d", len(out["misc"]))
	}
}

func TestRouter_SecurityRule_RoutesCorrectly(t *testing.T) {
	r := NewRouter(
		WithRouteRule("security", "sec-channel"),
	)
	out := r.Route(sampleRouterEntries)
	// libA and libC carry "security"
	if len(out["sec-channel"]) != 2 {
		t.Errorf("expected 2 entries in sec-channel, got %d", len(out["sec-channel"]))
	}
}

func TestRouter_MultipleRules_EntryInMultipleChannels(t *testing.T) {
	r := NewRouter(
		WithRouteRule("security", "sec-channel"),
		WithRouteRule("breaking", "breaking-channel"),
	)
	out := r.Route(sampleRouterEntries)
	// libC has both "security" and "breaking"
	if len(out["breaking-channel"]) != 1 {
		t.Errorf("expected 1 entry in breaking-channel, got %d", len(out["breaking-channel"]))
	}
	if len(out["sec-channel"]) != 2 {
		t.Errorf("expected 2 entries in sec-channel, got %d", len(out["sec-channel"]))
	}
}

func TestRouter_NoMatch_GoesToFallback(t *testing.T) {
	r := NewRouter(
		WithRouteRule("security", "sec-channel"),
	)
	out := r.Route(sampleRouterEntries)
	// libB (feature) and libD (no labels) fall through
	if len(out["default"]) != 2 {
		t.Errorf("expected 2 entries in default, got %d", len(out["default"]))
	}
}

func TestRouter_EmptyInput_ReturnsEmptyMap(t *testing.T) {
	r := NewRouter(WithRouteRule("security", "sec-channel"))
	out := r.Route([]Entry{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d channels", len(out))
	}
}
