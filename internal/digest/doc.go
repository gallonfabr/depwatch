// Package digest provides utilities for building and formatting dependency
// update digests from parsed changelog entries.
//
// A Builder collects changelog.Entry values keyed by dependency name and
// assembles them into a Digest, which can be rendered as plain text for
// delivery via Slack or email notifiers.
//
// Typical usage:
//
//	builder := digest.NewBuilder()
//	d := builder.Build(updates)          // updates: map[string][]changelog.Entry
//	body := d.FormatText()               // plain-text representation
package digest
