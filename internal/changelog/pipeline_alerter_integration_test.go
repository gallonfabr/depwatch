package changelog

import (
	"testing"
)

func TestPipeline_Alerter_SecurityEntryTriggersAlert(t *testing.T) {
	entries := []Entry{
		{Dependency: "dep-a", Version: "1.2.0", Tags: []string{"security"}},
		{Dependency: "dep-b", Version: "2.0.0", Tags: []string{"feature"}},
	}

	p := NewPipeline()
	out := p.Run(entries)

	alerter := NewAlerter(WithAlertRule("security", "high"))
	alerts := alerter.Evaluate(out)

	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert after pipeline, got %d", len(alerts))
	}
	if alerts[0].Entry.Dependency != "dep-a" {
		t.Errorf("expected alert for dep-a, got %q", alerts[0].Entry.Dependency)
	}
}

func TestPipeline_Alerter_NoTaggedEntries_NoAlerts(t *testing.T) {
	entries := []Entry{
		{Dependency: "dep-c", Version: "0.9.0", Tags: []string{"docs"}},
	}

	p := NewPipeline()
	out := p.Run(entries)

	alerter := NewAlerter(WithAlertRule("security", "high"))
	alerts := alerter.Evaluate(out)

	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestPipeline_Alerter_EmptyInput_DoesNotPanic(t *testing.T) {
	p := NewPipeline()
	out := p.Run([]Entry{})

	alerter := NewAlerter(WithAlertRule("breaking", "critical"))
	alerts := alerter.Evaluate(out)

	if alerts == nil {
		return // nil slice is acceptable
	}
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts for empty input, got %d", len(alerts))
	}
}
