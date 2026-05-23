package changelog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GitHubRelease represents a single GitHub release entry.
type GitHubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
	HTMLURL     string    `json:"html_url"`
}

// GitHubFetcher fetches changelog entries from GitHub Releases API.
type GitHubFetcher struct {
	client  *http.Client
	baseURL string
}

// NewGitHubFetcher creates a new GitHubFetcher with the given HTTP client.
// If client is nil, a default client with a 10-second timeout is used.
func NewGitHubFetcher(client *http.Client) *GitHubFetcher {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &GitHubFetcher{
		client:  client,
		baseURL: "https://api.github.com",
	}
}

// FetchReleases fetches the latest releases for the given owner/repo.
// It returns up to maxReleases entries.
func (g *GitHubFetcher) FetchReleases(owner, repo string, maxReleases int) ([]GitHubRelease, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo must not be empty")
	}
	if maxReleases <= 0 {
		maxReleases = 5
	}

	url := fmt.Sprintf("%s/repos/%s/%s/releases?per_page=%d", g.baseURL, owner, repo, maxReleases)
	resp, err := g.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching releases for %s/%s: %w", owner, repo, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s/%s", resp.StatusCode, owner, repo)
	}

	var releases []GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("decoding releases for %s/%s: %w", owner, repo, err)
	}
	return releases, nil
}
