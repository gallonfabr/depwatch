// Package changelog provides the Matcher type, which filters changelog
// entries based on one or more regular-expression rules.
//
// Each rule targets a specific field ("version" or "body") and is compiled
// once at construction time. Apply returns the union of all matching entries,
// preserving the original order. When no rules are registered, Apply is a
// no-op and returns all entries unchanged.
package changelog
