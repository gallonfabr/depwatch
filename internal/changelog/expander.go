package changelog

import "strings"

// Expander resolves shorthand version references (e.g. "v1" → "v1.0.0") in
// Entry.Version fields so that downstream semver comparisons work reliably.
//
// Segments missing from the original string are filled with zeroes:
//   - "v1"     → "v1.0.0"
//   - "v1.2"   → "v1.2.0"
//   - "v1.2.3" → "v1.2.3" (unchanged)
type Expander struct{}

// NewExpander returns a new Expander.
func NewExpander() *Expander { return &Expander{} }

// Apply normalises the Version field of every entry in the slice.
// The original slice is mutated in-place and also returned for chaining.
func (e *Expander) Apply(entries []Entry) []Entry {
	for i := range entries {
		entries[i].Version = expand(entries[i].Version)
	}
	return entries
}

// expand pads a version string to three dot-separated numeric segments.
func expand(v string) string {
	if v == "" {
		return v
	}

	prefix := ""
	raw := v
	if strings.HasPrefix(v, "v") || strings.HasPrefix(v, "V") {
		prefix = string(v[0])
		raw = v[1:]
	}

	// Preserve pre-release / build-metadata suffixes.
	core := raw
	suffix := ""
	if idx := strings.IndexAny(raw, "-+"); idx != -1 {
		core = raw[:idx]
		suffix = raw[idx:]
	}

	parts := strings.Split(core, ".")
	for len(parts) < 3 {
		parts = append(parts, "0")
	}

	return prefix + strings.Join(parts[:3], ".") + suffix
}
