package changelog

import (
	"testing"
	"time"
)

var sampleSplitterEntries = []Entry{
	{Dependency: "cobra", Version: "1.0.0", Date: time.Now()},
	{Dependency: "viper", Version: "2.0.0", Date: time.Now()},
	{Dependency: "cobra", Version: "1.1.0", Date: time.Now()},
	{Dependency: "viper", Version: "2.1.0", Date: time.Now()},
	{Dependency: "zap", Version: "3.0.0", Date: time.Now()},
}

func TestNewSplitter_NotNil(t *testing.T) {
	if NewSplitter() == nil {
		t.Fatal("expected non-nil Splitter")
	}
}

func TestSplitter_Split_GroupsByDependency(t *testing.T) {
	s := NewSplitter()
	buckets := s.Split(sampleSplitterEntries)
	if len(buckets["cobra"]) != 2 {
		t.Fatalf("expected 2 cobra entries, got %d", len(buckets["cobra"]))
	}
	if len(buckets["viper"]) != 2 {
		t.Fatalf("expected 2 viper entries, got %d", len(buckets["viper"]))
	}
	if len(buckets["zap"]) != 1 {
		t.Fatalf("expected 1 zap entry, got %d", len(buckets["zap"]))
	}
}

func TestSplitter_Split_EmptyInput(t *testing.T) {
	s := NewSplitter()
	buckets := s.Split(nil)
	if len(buckets) != 0 {
		t.Fatalf("expected empty map, got %d keys", len(buckets))
	}
}

func TestSplitter_Split_PreservesOrder(t *testing.T) {
	s := NewSplitter()
	buckets := s.Split(sampleSplitterEntries)
	cobra := buckets["cobra"]
	if cobra[0].Version != "1.0.0" || cobra[1].Version != "1.1.0" {
		t.Fatal("cobra entries are not in original order")
	}
}

func TestSplitter_Keys_FirstSeenOrder(t *testing.T) {
	s := NewSplitter()
	keys := s.Keys(sampleSplitterEntries)
	expected := []string{"cobra", "viper", "zap"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}
	for i, k := range expected {
		if keys[i] != k {
			t.Fatalf("key[%d]: expected %q, got %q", i, k, keys[i])
		}
	}
}

func TestSplitter_Keys_EmptyInput(t *testing.T) {
	s := NewSplitter()
	keys := s.Keys(nil)
	if len(keys) != 0 {
		t.Fatalf("expected no keys, got %d", len(keys))
	}
}

func TestSplitter_Split_EmptyDependency_CapturedUnderEmptyKey(t *testing.T) {
	s := NewSplitter()
	entries := []Entry{
		{Dependency: "", Version: "0.1.0"},
		{Dependency: "cobra", Version: "1.0.0"},
	}
	buckets := s.Split(entries)
	if len(buckets[""] ) != 1 {
		t.Fatalf("expected 1 entry under empty key, got %d", len(buckets[""]))
	}
}
