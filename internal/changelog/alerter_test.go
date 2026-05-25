package changelog

import (
	"testing"
)

var sampleAlerterEntries = []Entry{
	{Dependency: "libA", Version: "1.0.0", Tags: []string{"security", "breaking"}},
	{Dependency: "libB", Version: "2.0.0", Tags: []string{"feature"}},
	{Dependency: "libC", Version: "3.0.0", Tags: []string{}},
	{Dependency: "libD", Version: "4.0.0", Tags: []string{"bugfix", "security"}},
}

func TestNewAlerter_NotNil(t *testing.T) {
	a := NewAlerter()
	if a == nil {
		t.Fatal("expected non-nil Alerter")
	}
}

func TestAlerter_NoRules_NoAlerts(t *testing.T) {
	a := NewAlerter()
	alerts := a.Evaluate(sampleAlerterEntries)
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestAlerter_SecurityRule_MatchesTwoEntries(t *testing.T) {
	a := NewAlerter(WithAlertRule("security", "high"))
	alerts := a.Evaluate(sampleAlerterEntries)
	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}
}

func TestAlerter_Alert_SeveritySet(t *testing.T) {
	a := NewAlerter(WithAlertRule("security", "critical"))
	alerts := a.Evaluate(sampleAlerterEntries)
	for _, al := range alerts {
		if al.Severity != "critical" {
			t.Errorf("expected severity 'critical', got %q", al.Severity)
		}
	}
}

func TestAlerter_Alert_ReasonContainsTag(t *testing.T) {
	a := NewAlerter(WithAlertRule("breaking", "high"))
	alerts := a.Evaluate(sampleAlerterEntries)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Reason != "tag:breaking" {
		t.Errorf("unexpected reason: %q", alerts[0].Reason)
	}
}

func TestAlerter_MultipleRules_EachFires(t *testing.T) {
	a := NewAlerter(
		WithAlertRule("security", "high"),
		WithAlertRule("feature", "low"),
	)
	alerts := a.Evaluate(sampleAlerterEntries)
	if len(alerts) != 3 {
		t.Fatalf("expected 3 alerts, got %d", len(alerts))
	}
}

func TestAlerter_CaseInsensitiveTag(t *testing.T) {
	entries := []Entry{
		{Dependency: "libX", Tags: []string{"SECURITY"}},
	}
	a := NewAlerter(WithAlertRule("security", "high"))
	alerts := a.Evaluate(entries)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert for case-insensitive match, got %d", len(alerts))
	}
}

func TestAlerter_EmptyEntries_NoAlerts(t *testing.T) {
	a := NewAlerter(WithAlertRule("security", "high"))
	alerts := a.Evaluate([]Entry{})
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(alerts))
	}
}
