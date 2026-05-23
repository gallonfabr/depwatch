// Package changelog provides functionality for fetching and parsing
// dependency changelogs from HTTP URLs and GitHub releases.
//
// It includes:
//   - HTTPFetcher: fetches raw changelog content over HTTP.
//   - GitHubFetcher: fetches release notes via the GitHub Releases API.
//   - Parser: parses Markdown changelogs into structured Entry values.
//   - Cache: an in-memory TTL cache to avoid redundant network requests.
package changelog
