package changelog

import (
	"testing"
	"time"
)

var sampleTrendingEntries = []Entry{
	{Dependency: "react", Score: 10, Date: time.Now()},
	{Dependency: "react", Score: 5, Date: time.Now()},
	{Dependency: "lodash", Score: 20, Date: time.Now()},
	{Dependency: "axios", Score: 3, Date: time.Now()},
	{Dependency: "axios", Score: 3, Date: time.Now()},
	{Dependency: "axios", Score: 3, Date: time.Now()},
}

func TestNewTrending_NotNil(t *testing.T) {
	tr := NewTrending(5)
	if tr == nil {
		t.Fatal("expected non-nil Trending")
	}
}

func TestTrending_Analyse_OrderByScore(t *testing.T) {
	tr := NewTrending(0)
	result := tr.Analyse(sampleTrendingEntries)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[0].Dependency != "lodash" {
		t.Errorf("expected lodash first, got %s", result[0].Dependency)
	}
}

func TestTrending_Analyse_TopN(t *testing.T) {
	tr := NewTrending(2)
	result := tr.Analyse(sampleTrendingEntries)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestTrending_Analyse_Empty(t *testing.T) {
	tr := NewTrending(5)
	result := tr.Analyse(nil)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestTrending_Analyse_SkipsEmptyDependency(t *testing.T) {
	entries := []Entry{
		{Dependency: "", Score: 100},
		{Dependency: "react", Score: 5},
	}
	tr := NewTrending(0)
	result := tr.Analyse(entries)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Dependency != "react" {
		t.Errorf("unexpected dependency: %s", result[0].Dependency)
	}
}

func TestTrending_Analyse_TieBreakByCount(t *testing.T) {
	entries := []Entry{
		{Dependency: "alpha", Score: 5},
		{Dependency: "beta", Score: 5},
		{Dependency: "beta", Score: 0},
	}
	tr := NewTrending(0)
	result := tr.Analyse(entries)
	if result[0].Dependency != "beta" {
		t.Errorf("expected beta first due to higher count, got %s", result[0].Dependency)
	}
}

func TestTrending_Analyse_AggregatesCount(t *testing.T) {
	tr := NewTrending(0)
	result := tr.Analyse(sampleTrendingEntries)
	for _, r := range result {
		if r.Dependency == "axios" && r.Count != 3 {
			t.Errorf("expected axios count 3, got %d", r.Count)
		}
	}
}
