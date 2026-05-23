package changelog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewGitHubFetcher_DefaultClient(t *testing.T) {
	f := NewGitHubFetcher(nil)
	if f == nil {
		t.Fatal("expected non-nil fetcher")
	}
	if f.client == nil {
		t.Error("expected default HTTP client to be set")
	}
}

func TestNewGitHubFetcher_CustomClient(t *testing.T) {
	custom := &http.Client{Timeout: 5 * time.Second}
	f := NewGitHubFetcher(custom)
	if f.client != custom {
		t.Error("expected custom client to be used")
	}
}

func TestFetchReleases_Success(t *testing.T) {
	releases := []GitHubRelease{
		{TagName: "v1.2.0", Name: "Release 1.2.0", Body: "Bug fixes", PublishedAt: time.Now()},
		{TagName: "v1.1.0", Name: "Release 1.1.0", Body: "New features", PublishedAt: time.Now()},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(releases)
	}))
	defer server.Close()

	f := NewGitHubFetcher(server.Client())
	f.baseURL = server.URL

	got, err := f.FetchReleases("owner", "repo", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 releases, got %d", len(got))
	}
	if got[0].TagName != "v1.2.0" {
		t.Errorf("expected tag v1.2.0, got %s", got[0].TagName)
	}
}

func TestFetchReleases_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	f := NewGitHubFetcher(server.Client())
	f.baseURL = server.URL

	_, err := f.FetchReleases("owner", "repo", 5)
	if err == nil {
		t.Error("expected error for non-OK status")
	}
}

func TestFetchReleases_EmptyOwnerOrRepo(t *testing.T) {
	f := NewGitHubFetcher(nil)

	_, err := f.FetchReleases("", "repo", 5)
	if err == nil {
		t.Error("expected error for empty owner")
	}

	_, err = f.FetchReleases("owner", "", 5)
	if err == nil {
		t.Error("expected error for empty repo")
	}
}

func TestFetchReleases_DefaultMaxReleases(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		perPage := r.URL.Query().Get("per_page")
		if perPage != "5" {
			t.Errorf("expected per_page=5, got %s", perPage)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]GitHubRelease{})
	}))
	defer server.Close()

	f := NewGitHubFetcher(server.Client())
	f.baseURL = server.URL

	_, err := f.FetchReleases("owner", "repo", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
