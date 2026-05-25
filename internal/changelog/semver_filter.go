package changelog

import "fmt"

// SemVerFilter drops entries whose parsed version falls outside a
// caller-supplied minimum/maximum range.  Either bound may be left as
// the zero value to indicate "no bound".
type SemVerFilter struct {
	min *Version
	max *Version
}

// SemVerFilterOption configures a SemVerFilter.
type SemVerFilterOption func(*SemVerFilter)

// WithMinVersion sets the inclusive lower bound.
func WithMinVersion(v string) SemVerFilterOption {
	return func(f *SemVerFilter) {
		parsed, err := ParseVersion(v)
		if err == nil {
			f.min = &parsed
		}
	}
}

// WithMaxVersion sets the inclusive upper bound.
func WithMaxVersion(v string) SemVerFilterOption {
	return func(f *SemVerFilter) {
		parsed, err := ParseVersion(v)
		if err == nil {
			f.max = &parsed
		}
	}
}

// NewSemVerFilter constructs a SemVerFilter with the given options.
func NewSemVerFilter(opts ...SemVerFilterOption) *SemVerFilter {
	f := &SemVerFilter{}
	for _, o := range opts {
		o(f)
	}
	return f
}

// Apply removes entries that fall outside the configured version range.
// Entries whose version cannot be parsed are kept as-is.
func (f *SemVerFilter) Apply(entries []Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		v, err := ParseVersion(e.Version)
		if err != nil {
			// unparseable version — pass through
			out = append(out, e)
			continue
		}
		if f.min != nil && versionLess(v, *f.min) {
			continue
		}
		if f.max != nil && versionLess(*f.max, v) {
			continue
		}
		out = append(out, e)
	}
	return out
}

// versionLess returns true when a < b using major.minor.patch ordering.
func versionLess(a, b Version) bool {
	switch {
	case a.Major != b.Major:
		return a.Major < b.Major
	case a.Minor != b.Minor:
		return a.Minor < b.Minor
	case a.Patch != b.Patch:
		return a.Patch < b.Patch
	default:
		return false
	}
}

// String returns a human-readable description of the active range.
func (f *SemVerFilter) String() string {
	min, max := "*", "*"
	if f.min != nil {
		min = f.min.String()
	}
	if f.max != nil {
		max = f.max.String()
	}
	return fmt.Sprintf("SemVerFilter[%s, %s]", min, max)
}
