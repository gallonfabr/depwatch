package changelog

import (
	"testing"
)

var sampleLabelerEntries = []Entry{
	{Version: "1.0.0", Body: "Fixed a critical security vulnerability CVE-2024-1234"},
	{Version: "1.1.0", Body: "Added new feature for exporting data"},
	{Version: "1.2.0", Body: "Breaking change: removed legacy API endpoint"},
	{Version: "1.3.0", Body: "Bug fix for nil pointer dereference"},
	{Version: "1.4.0", Body: "Updated internal dependencies"},
}

func TestNewLabeler_NotNil(t *testing.T) {
	l := NewLabeler()
	if l == nil {
		t.Fatal("expected non-nil Labeler")
	}
}

func TestLabeler_Apply_SecurityLabel(t *testing.T) {
	l := NewLabeler()
	out := l.Apply(sampleLabelerEntries[:1])
	if out[0].Label != LabelSecurity {
		t.Errorf("expected %q, got %q", LabelSecurity, out[0].Label)
	}
}

func TestLabeler_Apply_FeatureLabel(t *testing.T) {
	l := NewLabeler()
	out := l.Apply(sampleLabelerEntries[1:2])
	if out[0].Label != LabelFeature {
		t.Errorf("expected %q, got %q", LabelFeature, out[0].Label)
	}
}

func TestLabeler_Apply_BreakingLabel(t *testing.T) {
	l := NewLabeler()
	out := l.Apply(sampleLabelerEntries[2:3])
	if out[0].Label != LabelBreaking {
		t.Errorf("expected %q, got %q", LabelBreaking, out[0].Label)
	}
}

func TestLabeler_Apply_BugfixLabel(t *testing.T) {
	l := NewLabeler()
	out := l.Apply(sampleLabelerEntries[3:4])
	if out[0].Label != LabelBugfix {
		t.Errorf("expected %q, got %q", LabelBugfix, out[0].Label)
	}
}

func TestLabeler_Apply_UnknownLabel(t *testing.T) {
	l := NewLabeler()
	out := l.Apply(sampleLabelerEntries[4:5])
	if out[0].Label != LabelUnknown {
		t.Errorf("expected %q, got %q", LabelUnknown, out[0].Label)
	}
}

func TestLabeler_Apply_DoesNotMutateInput(t *testing.T) {
	original := []Entry{{Version: "2.0.0", Body: "security patch applied"}}
	l := NewLabeler()
	_ = l.Apply(original)
	if original[0].Label != "" {
		t.Error("Apply must not mutate the input slice")
	}
}

func TestLabeler_Apply_SecurityTakesPriorityOverBugfix(t *testing.T) {
	entries := []Entry{
		{Version: "3.0.0", Body: "security fix for critical bug"},
	}
	l := NewLabeler()
	out := l.Apply(entries)
	if out[0].Label != LabelSecurity {
		t.Errorf("security should take priority, got %q", out[0].Label)
	}
}
