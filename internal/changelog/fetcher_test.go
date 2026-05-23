package changelog_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/changelog"
)

func TestHTTPFetcher_Fetch_Success(t *testing.T) {
	expectedContent := "## v1.2.3\n- fixed a bug\n"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(expectedContent))
	}))
	defer server.Close()

	f := changelog.NewHTTPFetcher(5 * time.Second)
	entry, err := f.Fetch("mylib", server.URL)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if entry.Dependency != "mylib" {
		t.Errorf("expected dependency %q, got %q", "mylib", entry.Dependency)
	}
	if entry.Content != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, entry.Content)
	}
	if entry.FetchedAt.IsZero() {
		t.Error("expected FetchedAt to be set")
	}
}

func TestHTTPFetcher_Fetch_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	f := changelog.NewHTTPFetcher(5 * time.Second)
	_, err := f.Fetch("mylib", server.URL)
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("expected error to mention status code, got: %v", err)
	}
}

func TestHTTPFetcher_Fetch_EmptyURL(t *testing.T) {
	f := changelog.NewHTTPFetcher(5 * time.Second)
	_, err := f.Fetch("mylib", "")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewHTTPFetcher_DefaultTimeout(t *testing.T) {
	f := changelog.NewHTTPFetcher(0)
	if f.Timeout != 10*time.Second {
		t.Errorf("expected default timeout 10s, got %v", f.Timeout)
	}
}
