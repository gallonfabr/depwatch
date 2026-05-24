package changelog

import (
	"testing"
	"time"
)

func sampleTransformEntries() []Entry {
	return []Entry{
		{Dependency: "alpha", Version: "1.0.0", Date: time.Now()},
		{Dependency: "beta", Version: "2.0.0", Date: time.Now()},
		{Dependency: "gamma", Version: "3.0.0", Date: time.Now()},
	}
}

func TestTransformFunc_ImplementsTransformer(t *testing.T) {
	var _ Transformer = TransformFunc(nil)
}

func TestTransformFunc_Apply(t *testing.T) {
	input := sampleTransformEntries()
	f := TransformFunc(func(entries []Entry) []Entry {
		return entries[:1]
	})
	out := f.Transform(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestNewChain_NotNil(t *testing.T) {
	if NewChain() == nil {
		t.Fatal("expected non-nil Chain")
	}
}

func TestChain_Empty_ReturnsUnchanged(t *testing.T) {
	c := NewChain()
	input := sampleTransformEntries()
	out := c.Transform(input)
	if len(out) != len(input) {
		t.Fatalf("expected %d entries, got %d", len(input), len(out))
	}
}

func TestChain_AppliesStepsInOrder(t *testing.T) {
	var order []string
	step1 := TransformFunc(func(e []Entry) []Entry { order = append(order, "first"); return e })
	step2 := TransformFunc(func(e []Entry) []Entry { order = append(order, "second"); return e })
	c := NewChain(step1, step2)
	c.Transform(sampleTransformEntries())
	if len(order) != 2 || order[0] != "first" || order[1] != "second" {
		t.Fatalf("unexpected order: %v", order)
	}
}

func TestNewLimitTransformer_NegativeBecomesZero(t *testing.T) {
	l := NewLimitTransformer(-5)
	out := l.Transform(sampleTransformEntries())
	if len(out) != 3 {
		t.Fatalf("expected 3 entries when max=0 (no-op), got %d", len(out))
	}
}

func TestLimitTransformer_Truncates(t *testing.T) {
	l := NewLimitTransformer(2)
	out := l.Transform(sampleTransformEntries())
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestLimitTransformer_UnderLimit_ReturnsAll(t *testing.T) {
	l := NewLimitTransformer(10)
	out := l.Transform(sampleTransformEntries())
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestChain_WithLimitTransformer(t *testing.T) {
	c := NewChain(NewLimitTransformer(1))
	out := c.Transform(sampleTransformEntries())
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}
