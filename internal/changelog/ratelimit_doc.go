// Package changelog provides utilities for fetching, parsing, and processing
// dependency changelogs from various sources.
//
// RateLimiter enforces a sliding-window rate limit per dependency key,
// preventing excessive requests to upstream changelog sources during a
// polling cycle. Use NewRateLimiter to configure the allowed call count
// and time window, then call Allow before each fetch operation.
package changelog
