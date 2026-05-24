package changelog

import (
	"testing"
	"time"
)

func makeEntry(dep, body string) Entry {
	v, _ := ParseVersion("v1.0.0")
	return Entry{Dependency: dep, Version: v, Date: time.Now(), Body: body}
}

func TestNewRanker_NilScorer_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil scorer")
		}
	}()
	NewRanker(nil)
}

func TestNewRanker_Valid(t *testing.T) {
	s := NewScorer()
	r := NewRanker(s)
	if r == nil {
		t.Fatal("expected non-nil Ranker")
	}
}

func TestRanker_Rank_HigherScoreFirst(t *testing.T) {
	s := NewScorer(WithKeywords("security"))
	r := NewRanker(s)

	entries := []Entry{
		makeEntry("a", "minor update"),
		makeEntry("b", "security vulnerability fixed"),
	}

	ranked := r.Rank(entries)
	if ranked[0].Dependency != "b" {
		t.Fatalf("expected 'b' first, got %s", ranked[0].Dependency)
	}
}

func TestRanker_Rank_DoesNotMutateInput(t *testing.T) {
	s := NewScorer(WithKeywords("fix"))
	r := NewRanker(s)

	orig := []Entry{
		makeEntry("a", "minor update"),
		makeEntry("b", "fix crash"),
	}
	origFirst := orig[0].Dependency

	r.Rank(orig)

	if orig[0].Dependency != origFirst {
		t.Fatal("Rank must not mutate input slice")
	}
}

func TestRanker_TopN_LimitsResults(t *testing.T) {
	s := NewScorer(WithKeywords("fix"))
	r := NewRanker(s)

	entries := []Entry{
		makeEntry("a", "fix one"),
		makeEntry("b", "fix two"),
		makeEntry("c", "fix three"),
	}

	top := r.TopN(entries, 2)
	if len(top) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(top))
	}
}

func TestRanker_TopN_ZeroReturnsAll(t *testing.T) {
	s := NewScorer()
	r := NewRanker(s)
	entries := []Entry{makeEntry("a", "x"), makeEntry("b", "y")}
	if got := r.TopN(entries, 0); len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
