// Package changelog provides utilities for fetching, parsing, and processing
// dependency changelogs from various sources such as HTTP endpoints and the
// GitHub Releases API.
//
// The Throttle type in this file limits the rate at which individual
// dependency sources are polled, preventing excessive outbound requests when
// the scheduler interval is very short or when many dependencies are tracked.
package changelog
