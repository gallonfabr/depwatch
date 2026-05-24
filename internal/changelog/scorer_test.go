package changelog

import (
	"testing"
	"time"
)

func sampleScorerEntry(body string) Entry {
	v, _ := ParseVersion("v1.2.3")
	return Entry{
		Dependency: "lib",
		Version:    v,
		Date:       time.Now(),
		Body:       body,
	}
}

func TestNewScorer_NotNil(t *testing.T) {
	s := NewScorer()
	if s == nil {
		t.Fatal("expected non-nil Scorer")
	}
}

func TestScorer_NoKeywords_ZeroScore(t *testing.T) {
	s := NewScorer()
	e := sampleScorerEntry("breaking change in API")
	if got := s.Score(e); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestScorer_KeywordMatch_IncreasesScore(t *testing.T) {
	s := NewScorer(WithKeywords("breaking"))
	e := sampleScorerEntry("breaking change in API")
	if got := s.Score(e); got < 1 {
		t.Fatalf("expected score >= 1, got %d", got)
	}
}

func TestScorer_MultipleMatches_AccumulatesScore(t *testing.T) {
	s := NewScorer(WithKeywords("fix"))
	e := sampleScorerEntry("fix memory leak, fix race condition")
	if got := s.Score(e); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestScorer_SecurityKeyword_DoubleCounts(t *testing.T) {
	s := NewScorer(WithKeywords("security"))
	e := sampleScorerEntry("security fix for auth bypass")
	if got := s.Score(e); got != 2 {
		t.Fatalf("expected 2 (security weight), got %d", got)
	}
}

func TestScorer_CaseInsensitive(t *testing.T) {
	s := NewScorer(WithKeywords("BREAKING"))
	e := sampleScorerEntry("breaking change")
	if got := s.Score(e); got < 1 {
		t.Fatalf("expected score >= 1, got %d", got)
	}
}

func TestScorer_NoMatch_ZeroScore(t *testing.T) {
	s := NewScorer(WithKeywords("deprecat"))
	e := sampleScorerEntry("minor style fixes")
	if got := s.Score(e); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
