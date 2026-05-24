package changelog

import (
	"testing"
)

func sampleGrouperEntries() []Entry {
	return []Entry{
		{Version: "1.0.0", Labels: []string{"security"}},
		{Version: "1.1.0", Labels: []string{"feature"}},
		{Version: "1.2.0", Labels: []string{"bugfix"}},
		{Version: "1.3.0", Labels: []string{"security"}},
		{Version: "1.4.0"},
	}
}

func TestNewGrouper_NotNil(t *testing.T) {
	g := NewGrouper()
	if g == nil {
		t.Fatal("expected non-nil Grouper")
	}
}

func TestGrouper_DefaultFallback(t *testing.T) {
	g := NewGrouper()
	groups := g.Apply([]Entry{{Version: "1.0.0"}})
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Label != "other" {
		t.Errorf("expected fallback label 'other', got %q", groups[0].Label)
	}
}

func TestGrouper_CustomFallback(t *testing.T) {
	g := NewGrouper(WithFallbackLabel("misc"))
	groups := g.Apply([]Entry{{Version: "1.0.0"}})
	if groups[0].Label != "misc" {
		t.Errorf("expected fallback label 'misc', got %q", groups[0].Label)
	}
}

func TestGrouper_GroupsByFirstLabel(t *testing.T) {
	g := NewGrouper()
	entries := sampleGrouperEntries()
	groups := g.Apply(entries)

	total := 0
	for _, grp := range groups {
		total += len(grp.Entries)
	}
	if total != len(entries) {
		t.Errorf("expected %d total entries across groups, got %d", len(entries), total)
	}
}

func TestGrouper_SecurityGroupHasTwoEntries(t *testing.T) {
	g := NewGrouper()
	groups := g.Apply(sampleGrouperEntries())

	for _, grp := range groups {
		if grp.Label == "security" && len(grp.Entries) != 2 {
			t.Errorf("expected 2 security entries, got %d", len(grp.Entries))
		}
	}
}

func TestGrouper_OrderRespected(t *testing.T) {
	g := NewGrouper(WithGroupOrder("security", "feature", "bugfix"))
	groups := g.Apply(sampleGrouperEntries())

	if len(groups) < 3 {
		t.Fatalf("expected at least 3 groups, got %d", len(groups))
	}
	expected := []string{"security", "feature", "bugfix"}
	for i, label := range expected {
		if groups[i].Label != label {
			t.Errorf("position %d: expected %q, got %q", i, label, groups[i].Label)
		}
	}
}

func TestGrouper_EmptyInput(t *testing.T) {
	g := NewGrouper()
	groups := g.Apply([]Entry{})
	if len(groups) != 0 {
		t.Errorf("expected 0 groups for empty input, got %d", len(groups))
	}
}

func TestGrouper_WithGroupOrder_EmptyLabelIgnored(t *testing.T) {
	g := NewGrouper(WithFallbackLabel(""))
	if g.fallback != "other" {
		t.Errorf("empty fallback label should not override default, got %q", g.fallback)
	}
}
