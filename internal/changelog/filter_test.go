package changelog

import (
	"testing"
	"time"
)

var baseTime = time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

func sampleEntries() []Entry {
	return []Entry{
		{Version: "v1.3.0", Date: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC), Body: "new feature"},
		{Version: "v1.2.0", Date: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), Body: "fix bug"},
		{Version: "v1.1.0", Date: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), Body: "initial"},
	}
}

func TestFilter_Apply_Since(t *testing.T) {
	f := NewFilter(baseTime, 0)
	got := f.Apply(sampleEntries())
	if len(got) != 2 {
		t.Fatalf("expected 2 entries after since filter, got %d", len(got))
	}
}

func TestFilter_Apply_Limit(t *testing.T) {
	f := NewFilter(time.Time{}, 1)
	got := f.Apply(sampleEntries())
	if len(got) != 1 {
		t.Fatalf("expected 1 entry with limit=1, got %d", len(got))
	}
	if got[0].Version != "v1.3.0" {
		t.Errorf("expected first entry v1.3.0, got %s", got[0].Version)
	}
}

func TestFilter_Apply_SinceAndLimit(t *testing.T) {
	f := NewFilter(baseTime, 1)
	got := f.Apply(sampleEntries())
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
}

func TestFilter_Apply_NoMatch(t *testing.T) {
	future := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	f := NewFilter(future, 0)
	got := f.Apply(sampleEntries())
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestFilter_Apply_ZeroSince_ReturnsAll(t *testing.T) {
	f := NewFilter(time.Time{}, 0)
	got := f.Apply(sampleEntries())
	if len(got) != 3 {
		t.Fatalf("expected all 3 entries, got %d", len(got))
	}
}

func TestFilter_Apply_EmptyEntries(t *testing.T) {
	f := NewFilter(baseTime, 5)
	got := f.Apply([]Entry{})
	if len(got) != 0 {
		t.Fatalf("expected 0 entries for empty input, got %d", len(got))
	}
}

func TestFilterNew_RemovesSeen(t *testing.T) {
	seen := map[string]bool{"v1.3.0": true, "v1.1.0": true}
	got := FilterNew(sampleEntries(), seen)
	if len(got) != 1 {
		t.Fatalf("expected 1 new entry, got %d", len(got))
	}
	if got[0].Version != "v1.2.0" {
		t.Errorf("expected v1.2.0, got %s", got[0].Version)
	}
}

func TestFilterNew_EmptySeen_ReturnsAll(t *testing.T) {
	got := FilterNew(sampleEntries(), map[string]bool{})
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestFilterNew_AllSeen_ReturnsNone(t *testing.T) {
	seen := map[string]bool{"v1.3.0": true, "v1.2.0": true, "v1.1.0": true}
	got := FilterNew(sampleEntries(), seen)
	if len(got) != 0 {
		t.Fatalf("expected 0 entries when all are seen, got %d", len(got))
	}
}
