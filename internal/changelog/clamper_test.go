package changelog

import (
	"testing"
)

func TestNewClamper_NotNil(t *testing.T) {
	c := NewClamper()
	if c == nil {
		t.Fatal("expected non-nil Clamper")
	}
}

func TestNewClamper_DefaultRange(t *testing.T) {
	c := NewClamper()
	if c.min != 0 || c.max != 100 {
		t.Fatalf("expected default range [0,100], got [%v,%v]", c.min, c.max)
	}
}

func TestNewClamper_CustomRange(t *testing.T) {
	c := NewClamper(WithClampMin(10), WithClampMax(50))
	if c.min != 10 || c.max != 50 {
		t.Fatalf("expected [10,50], got [%v,%v]", c.min, c.max)
	}
}

func TestNewClamper_InvertedRange_Swapped(t *testing.T) {
	c := NewClamper(WithClampMin(80), WithClampMax(20))
	if c.min != 20 || c.max != 80 {
		t.Fatalf("expected swapped range [20,80], got [%v,%v]", c.min, c.max)
	}
}

func TestClamper_Apply_BelowMin_ClampsUp(t *testing.T) {
	c := NewClamper(WithClampMin(5), WithClampMax(50))
	entries := []Entry{{Score: -3}, {Score: 0}}
	result := c.Apply(entries)
	for _, e := range result {
		if e.Score < 5 {
			t.Fatalf("expected score >= 5, got %v", e.Score)
		}
	}
}

func TestClamper_Apply_AboveMax_ClampsDown(t *testing.T) {
	c := NewClamper(WithClampMin(0), WithClampMax(10))
	entries := []Entry{{Score: 200}, {Score: 11}}
	result := c.Apply(entries)
	for _, e := range result {
		if e.Score > 10 {
			t.Fatalf("expected score <= 10, got %v", e.Score)
		}
	}
}

func TestClamper_Apply_WithinRange_Unchanged(t *testing.T) {
	c := NewClamper(WithClampMin(0), WithClampMax(100))
	entries := []Entry{{Score: 42}, {Score: 0}, {Score: 100}}
	result := c.Apply(entries)
	expected := []float64{42, 0, 100}
	for i, e := range result {
		if e.Score != expected[i] {
			t.Fatalf("entry %d: expected score %v, got %v", i, expected[i], e.Score)
		}
	}
}

func TestClamper_Apply_EmptyInput(t *testing.T) {
	c := NewClamper()
	result := c.Apply([]Entry{})
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(result))
	}
}
