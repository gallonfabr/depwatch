package changelog

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Entry represents a single changelog entry for a dependency.
type Entry struct {
	Dependency string
	Version    string
	Content    string
	FetchedAt  time.Time
}

// Fetcher defines the interface for fetching changelog content.
type Fetcher interface {
	Fetch(dep, url string) (*Entry, error)
}

// HTTPFetcher fetches changelogs over HTTP.
type HTTPFetcher struct {
	Client  *http.Client
	Timeout time.Duration
}

// NewHTTPFetcher creates a new HTTPFetcher with the given timeout.
func NewHTTPFetcher(timeout time.Duration) *HTTPFetcher {
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &HTTPFetcher{
		Client:  &http.Client{Timeout: timeout},
		Timeout: timeout,
	}
}

// Fetch retrieves the changelog content for a dependency from the given URL.
func (f *HTTPFetcher) Fetch(dep, url string) (*Entry, error) {
	if url == "" {
		return nil, fmt.Errorf("changelog URL for dependency %q is empty", dep)
	}

	resp, err := f.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching changelog for %q: %w", dep, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d fetching changelog for %q", resp.StatusCode, dep)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading changelog body for %q: %w", dep, err)
	}

	return &Entry{
		Dependency: dep,
		Content:    string(body),
		FetchedAt:  time.Now().UTC(),
	}, nil
}
