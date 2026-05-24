package changelog

import (
	"testing"
)

func TestNewClassifier_NotNil(t *testing.T) {
	c := NewClassifier()
	if c == nil {
		t.Fatal("expected non-nil Classifier")
	}
}

func TestClassifier_DefaultFallback(t *testing.T) {
	c := NewClassifier()
	entries := []Entry{{Version: "1.0.0", Labels: []string{}}}
	out := c.Apply(entries)
	if out[0].Category != "Other" {
		t.Errorf("expected 'Other', got %q", out[0].Category)
	}
}

func TestClassifier_CustomFallback(t *testing.T) {
	c := NewClassifier(WithClassifierFallback("Misc"))
	entries := []Entry{{Version: "1.0.0", Labels: []string{"unknown"}}}
	out := c.Apply(entries)
	if out[0].Category != "Misc" {
		t.Errorf("expected 'Misc', got %q", out[0].Category)
	}
}

func TestClassifier_SecurityLabel(t *testing.T) {
	c := NewClassifier()
	entries := []Entry{{Labels: []string{"security"}}}
	out := c.Apply(entries)
	if out[0].Category != "Security" {
		t.Errorf("expected 'Security', got %q", out[0].Category)
	}
}

func TestClassifier_FeatureLabel(t *testing.T) {
	c := NewClassifier()
	entries := []Entry{{Labels: []string{"feature"}}}
	out := c.Apply(entries)
	if out[0].Category != "Features" {
		t.Errorf("expected 'Features', got %q", out[0].Category)
	}
}

func TestClassifier_BugfixLabel(t *testing.T) {
	c := NewClassifier()
	entries := []Entry{{Labels: []string{"bugfix"}}}
	out := c.Apply(entries)
	if out[0].Category != "Bug Fixes" {
		t.Errorf("expected 'Bug Fixes', got %q", out[0].Category)
	}
}

func TestClassifier_FirstLabelWins(t *testing.T) {
	c := NewClassifier()
	entries := []Entry{{Labels: []string{"feature", "security"}}}
	out := c.Apply(entries)
	if out[0].Category != "Features" {
		t.Errorf("expected 'Features', got %q", out[0].Category)
	}
}

func TestClassifier_CustomMapping(t *testing.T) {
	c := NewClassifier(WithCategoryMapping("perf", "Performance"))
	entries := []Entry{{Labels: []string{"perf"}}}
	out := c.Apply(entries)
	if out[0].Category != "Performance" {
		t.Errorf("expected 'Performance', got %q", out[0].Category)
	}
}

func TestClassifier_DoesNotMutateInput(t *testing.T) {
	c := NewClassifier()
	orig := []Entry{{Labels: []string{"security"}, Category: ""}}
	_ = c.Apply(orig)
	if orig[0].Category != "" {
		t.Error("Apply must not mutate the original slice")
	}
}
