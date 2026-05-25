// Package changelog provides primitives for fetching, parsing, and
// processing dependency changelogs.
//
// Alerter evaluates a slice of changelog entries against a set of tag-based
// rules and returns Alert values for any matched entries. Each rule maps a
// tag name to a severity string (e.g. "high", "critical"). Only the first
// matching rule fires per entry to avoid duplicate alerts.
package changelog
