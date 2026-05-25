// Package changelog provides the Quota type for rate-limiting the number of
// changelog entries accepted per dependency within a rolling time window.
//
// Quota is useful when a single dependency publishes many releases in a short
// period and you want to avoid flooding Slack or email digests. Configure it
// with a maximum entry count and a window duration:
//
//	q := changelog.NewQuota(10, 24*time.Hour)
//	filtered := q.Apply(entries)
//
// Quota is safe for concurrent use. Call Reset to clear all recorded counts
// and start a fresh period.
package changelog
