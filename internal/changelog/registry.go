package changelog

import (
	"fmt"
	"strings"
)

// FetcherType represents the type of changelog fetcher to use.
type FetcherType string

const (
	FetcherTypeHTTP   FetcherType = "http"
	FetcherTypeGitHub FetcherType = "github"
)

// FetcherConfig holds the configuration needed to build a Fetcher.
type FetcherConfig struct {
	// Type determines which fetcher implementation to use.
	Type FetcherType

	// URL is used by the HTTP fetcher.
	URL string

	// Owner and Repo are used by the GitHub fetcher.
	Owner string
	Repo  string
}

// Registry creates Fetcher instances from FetcherConfig.
type Registry struct {
	httpFetcher   *HTTPFetcher
	githubFetcher *GitHubFetcher
}

// NewRegistry returns a Registry with default fetcher implementations.
func NewRegistry() *Registry {
	return &Registry{
		httpFetcher:   NewHTTPFetcher(0),
		githubFetcher: NewGitHubFetcher(nil),
	}
}

// Build returns a Fetcher for the given FetcherConfig.
func (r *Registry) Build(cfg FetcherConfig) (Fetcher, error) {
	switch FetcherType(strings.ToLower(string(cfg.Type))) {
	case FetcherTypeHTTP:
		if cfg.URL == "" {
			return nil, fmt.Errorf("changelog registry: http fetcher requires a URL")
		}
		return &boundHTTPFetcher{fetcher: r.httpFetcher, url: cfg.URL}, nil
	case FetcherTypeGitHub:
		if cfg.Owner == "" || cfg.Repo == "" {
			return nil, fmt.Errorf("changelog registry: github fetcher requires owner and repo")
		}
		return &boundGitHubFetcher{fetcher: r.githubFetcher, owner: cfg.Owner, repo: cfg.Repo}, nil
	default:
		return nil, fmt.Errorf("changelog registry: unknown fetcher type %q", cfg.Type)
	}
}

// boundHTTPFetcher wraps HTTPFetcher with a fixed URL.
type boundHTTPFetcher struct {
	fetcher *HTTPFetcher
	url     string
}

func (b *boundHTTPFetcher) Fetch() (string, error) {
	return b.fetcher.Fetch(b.url)
}

// boundGitHubFetcher wraps GitHubFetcher with a fixed owner/repo.
type boundGitHubFetcher struct {
	fetcher *GitHubFetcher
	owner   string
	repo    string
}

func (b *boundGitHubFetcher) Fetch() (string, error) {
	return b.fetcher.FetchReleases(b.owner, b.repo)
}
