// Package changelog provides utilities for fetching, parsing, and processing
// dependency changelogs from various sources.
//
// The retry sub-feature exposes Retry, a context-aware exponential-backoff
// helper used by fetchers to survive transient network failures. Wrap an error
// with Permanent to signal that no further attempts should be made.
package changelog
