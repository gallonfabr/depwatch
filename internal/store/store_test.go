package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/depwatch/internal/store"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func TestNew_CreatesEmptyStore(t *testing.T) {
	s, err := store.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.LastSeen("react"); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestSetLastSeen_PersistsValue(t *testing.T) {
	p := tempPath(t)
	s, _ := store.New(p)

	if err := s.SetLastSeen("lodash", "4.17.21"); err != nil {
		t.Fatalf("SetLastSeen error: %v", err)
	}
	if got := s.LastSeen("lodash"); got != "4.17.21" {
		t.Errorf("expected 4.17.21, got %q", got)
	}
}

func TestNew_LoadsExistingFile(t *testing.T) {
	p := tempPath(t)
	s1, _ := store.New(p)
	_ = s1.SetLastSeen("express", "4.18.2")

	s2, err := store.New(p)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	if got := s2.LastSeen("express"); got != "4.18.2" {
		t.Errorf("expected 4.18.2 after reload, got %q", got)
	}
}

func TestNew_InvalidJSON(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not-json"), 0o644)

	_, err := store.New(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestSetLastSeen_OverwritesValue(t *testing.T) {
	p := tempPath(t)
	s, _ := store.New(p)
	_ = s.SetLastSeen("axios", "1.0.0")
	_ = s.SetLastSeen("axios", "1.6.0")

	if got := s.LastSeen("axios"); got != "1.6.0" {
		t.Errorf("expected 1.6.0, got %q", got)
	}
}
