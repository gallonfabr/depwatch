package changelog

import (
	"testing"
	"time"
)

var sampleAggEntries = []Entry{
	{Dependency: "react", Labels: []string{"security", "feature"}, Highlighted: true, Date: time.Now()},
	{Dependency: "react", Labels: []string{"bugfix"}, Highlighted: false, Date: time.Now()},
	{Dependency: "lodash", Labels: []string{"security"}, Highlighted: true, Date: time.Now()},
	{Dependency: "webpack", Labels: []string{"feature"}, Highlighted: false, Date: time.Now()},
}

func TestNewAggregator_NotNil(t *testing.T) {
	if NewAggregator() == nil {
		t.Fatal("expected non-nil aggregator")
	}
}

func TestAggregate_TotalEntries(t *testing.T) {
	a := NewAggregator()
	stats := a.Aggregate(sampleAggEntries)
	if stats.TotalEntries != 4 {
		t.Fatalf("expected 4, got %d", stats.TotalEntries)
	}
}

func TestAggregate_ByDependency(t *testing.T) {
	a := NewAggregator()
	stats := a.Aggregate(sampleAggEntries)
	if stats.ByDependency["react"] != 2 {
		t.Fatalf("expected react=2, got %d", stats.ByDependency["react"])
	}
	if stats.ByDependency["lodash"] != 1 {
		t.Fatalf("expected lodash=1, got %d", stats.ByDependency["lodash"])
	}
}

func TestAggregate_ByLabel(t *testing.T) {
	a := NewAggregator()
	stats := a.Aggregate(sampleAggEntries)
	if stats.ByLabel["security"] != 2 {
		t.Fatalf("expected security=2, got %d", stats.ByLabel["security"])
	}
	if stats.ByLabel["feature"] != 2 {
		t.Fatalf("expected feature=2, got %d", stats.ByLabel["feature"])
	}
}

func TestAggregate_HighlightedCount(t *testing.T) {
	a := NewAggregator()
	stats := a.Aggregate(sampleAggEntries)
	if stats.HighlightedCount != 2 {
		t.Fatalf("expected 2, got %d", stats.HighlightedCount)
	}
}

func TestAggregate_EmptyEntries(t *testing.T) {
	a := NewAggregator()
	stats := a.Aggregate(nil)
	if stats.TotalEntries != 0 {
		t.Fatalf("expected 0, got %d", stats.TotalEntries)
	}
}

func TestTopDependencies_Order(t *testing.T) {
	a := NewAggregator()
	stats := a.Aggregate(sampleAggEntries)
	top := a.TopDependencies(stats, 0)
	if top[0] != "react" {
		t.Fatalf("expected react first, got %s", top[0])
	}
}

func TestTopDependencies_CapN(t *testing.T) {
	a := NewAggregator()
	stats := a.Aggregate(sampleAggEntries)
	top := a.TopDependencies(stats, 2)
	if len(top) != 2 {
		t.Fatalf("expected 2, got %d", len(top))
	}
}
