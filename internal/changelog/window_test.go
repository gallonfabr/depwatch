package changelog

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func sampleWindowEntries() []Entry {
	return []Entry{
		{Version: "v1.0.0", Date: epoch},
		{Version: "v1.1.0", Date: epoch.Add(24 * time.Hour)},
		{Version: "v1.2.0", Date: epoch.Add(48 * time.Hour)},
		{Version: "v1.3.0", Date: epoch.Add(72 * time.Hour)},
	}
}

func TestNewWindow_InvalidBounds(t *testing.T) {
	_, err := NewWindow(epoch.Add(time.Hour), epoch)
	if err == nil {
		t.Fatal("expected error for start >= end, got nil")
	}
}

func TestNewWindow_EqualBoundsInvalid(t *testing.T) {
	_, err := NewWindow(epoch, epoch)
	if err == nil {
		t.Fatal("expected error when start == end")
	}
}

func TestNewWindow_ValidBounds(t *testing.T) {
	w, err := NewWindow(epoch, epoch.Add(time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.IsZero() {
		t.Fatal("expected non-zero window")
	}
}

func TestNewWindow_ZeroBoundsAllowed(t *testing.T) {
	_, err := NewWindow(time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("unexpected error for zero bounds: %v", err)
	}
}

func TestWindow_Apply_StartBound(t *testing.T) {
	entries := sampleWindowEntries()
	w, _ := NewWindow(epoch.Add(24*time.Hour), time.Time{})
	got := w.Apply(entries)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestWindow_Apply_EndBound(t *testing.T) {
	entries := sampleWindowEntries()
	w, _ := NewWindow(time.Time{}, epoch.Add(48*time.Hour))
	got := w.Apply(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestWindow_Apply_BothBounds(t *testing.T) {
	entries := sampleWindowEntries()
	w, _ := NewWindow(epoch.Add(24*time.Hour), epoch.Add(72*time.Hour))
	got := w.Apply(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Version != "v1.1.0" || got[1].Version != "v1.2.0" {
		t.Fatalf("unexpected versions: %v", got)
	}
}

func TestWindow_Apply_Empty(t *testing.T) {
	w, _ := NewWindow(epoch, epoch.Add(time.Hour))
	got := w.Apply(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}

func TestWindow_Apply_NoMatch(t *testing.T) {
	entries := sampleWindowEntries()
	w, _ := NewWindow(epoch.Add(100*24*time.Hour), epoch.Add(200*24*time.Hour))
	got := w.Apply(entries)
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}
