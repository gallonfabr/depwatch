package changelog_test

import (
	"testing"

	"github.com/depwatch/internal/changelog"
)

func TestNewRegistry_NotNil(t *testing.T) {
	r := changelog.NewRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestRegistry_Build_HTTPFetcher_MissingURL(t *testing.T) {
	r := changelog.NewRegistry()
	_, err := r.Build(changelog.FetcherConfig{
		Type: changelog.FetcherTypeHTTP,
		URL:  "",
	})
	if err == nil {
		t.Fatal("expected error for missing URL, got nil")
	}
}

func TestRegistry_Build_HTTPFetcher_Valid(t *testing.T) {
	r := changelog.NewRegistry()
	f, err := r.Build(changelog.FetcherConfig{
		Type: changelog.FetcherTypeHTTP,
		URL:  "https://example.com/CHANGELOG.md",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil fetcher")
	}
}

func TestRegistry_Build_GitHubFetcher_MissingOwner(t *testing.T) {
	r := changelog.NewRegistry()
	_, err := r.Build(changelog.FetcherConfig{
		Type:  changelog.FetcherTypeGitHub,
		Owner: "",
		Repo:  "myrepo",
	})
	if err == nil {
		t.Fatal("expected error for missing owner, got nil")
	}
}

func TestRegistry_Build_GitHubFetcher_MissingRepo(t *testing.T) {
	r := changelog.NewRegistry()
	_, err := r.Build(changelog.FetcherConfig{
		Type:  changelog.FetcherTypeGitHub,
		Owner: "myowner",
		Repo:  "",
	})
	if err == nil {
		t.Fatal("expected error for missing repo, got nil")
	}
}

func TestRegistry_Build_GitHubFetcher_Valid(t *testing.T) {
	r := changelog.NewRegistry()
	f, err := r.Build(changelog.FetcherConfig{
		Type:  changelog.FetcherTypeGitHub,
		Owner: "myowner",
		Repo:  "myrepo",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil fetcher")
	}
}

func TestRegistry_Build_UnknownType(t *testing.T) {
	r := changelog.NewRegistry()
	_, err := r.Build(changelog.FetcherConfig{
		Type: "rss",
	})
	if err == nil {
		t.Fatal("expected error for unknown type, got nil")
	}
}

func TestRegistry_Build_CaseInsensitiveType(t *testing.T) {
	r := changelog.NewRegistry()
	f, err := r.Build(changelog.FetcherConfig{
		Type: "HTTP",
		URL:  "https://example.com/CHANGELOG.md",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil fetcher")
	}
}
